package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"technique/helper"
	"technique/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

//connection avec la base de données
var collection = helper.ConnectDB()

type M = primitive.M

func CheckLogin(c *gin.Context) {
	var login models.Login
	var user models.Login
	//je récupére le payload qui contient les données en json
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))
	//on recupère les données, la conversion en STRUCT
	err = json.Unmarshal(jsonData, &login)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("ID = %s, PASSWORD=%s", login.ID, login.Password)
	//il faut crypter le MDP pourqu'on puisse le comparer avec celui stocké dans la base de données (sinon on peut utiliser la fonction CompareHashWithPassword)
	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	login.Password = string(passwordHashed)
	//je cherche l'utilisateur dans la bd avec son id
	err = collection.FindOne(context.TODO(), bson.M{"id": login.ID}).Decode(&user)
	//si l'erreur est != nil ===> utilisateur non existant
	if err != nil {
		fmt.Println("Utilisateur inexistant")
		c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Can't find the user"})
	} else {
		if login.Password == user.Password {
			c.JSON(http.StatusOK, gin.H{"status": "sucess", "isCorrect": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "isCorrect": false, "message": "Mot de passe incorrect"})
		}

	}
}
func GetUsers(c *gin.Context) {
	//la création d'un tableau d'USERS
	var users []models.User
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, elem)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(context.TODO())

	//la sérialisation en JSON
	c.JSON(http.StatusOK, gin.H{"status": "success", "users": &users})

}

//la fonction pour récupérer un utilisateur avec son id
func GetUserById(c *gin.Context) {
	//l'id de l'utilisateur
	ID := c.Param("id")
	var user models.User
	//je cherche l'utilisateur dans la bd avec son id
	err := collection.FindOne(context.TODO(), bson.M{"id": ID}).Decode(&user)
	//si l'erreur est != nil ===> utilisateur non existant
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Can't find the user"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "sucess", "user": &user})
	}
}

//la fonction pour ajouter des utiisateurs dans la base de données.
func AddUsers(c *gin.Context) {
	//la création d'un waitgroup pour gérer les goroutines
	var waitGroup sync.WaitGroup
	//on recupère le chemin vers fichier json
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var users []models.User
	unmarshalErr := json.Unmarshal(jsonData, &users)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}
	fmt.Println("Number of users to add : " + string(len(users)))
	waitGroup.Add(len(users))
	for i := 0; i < len(users); i++ {
		users[i].Password = hashingPassword(users[i].Password)
		go addUser(&users[i], &waitGroup)
	}
	// On attends la find d'insertion de tous les utilisateurs.
	waitGroup.Wait()
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "All the users have been inserted!!!"})
}
func addUser(user *models.User, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
		fmt.Println("InsertOne() API result:", result)
		//la création du fichier (le nom sera l'ID de l'utilisateur), je vais stocker data dedans, on donne 755 comme permission
		err := ioutil.WriteFile("./files/"+user.ID, []byte(user.Data), 0755)
		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		} else {
			fmt.Printf("Fichier " + user.ID + " a ete creee avec succés.")
		}
		fmt.Println("CREATING FILE")
	}
}

//la fonction pour hasher le mot de passe
func hashingPassword(mdp string) string {
	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(mdp), bcrypt.DefaultCost)
	return string(passwordHashed)
}

//supprimer un utilisateur
func DeleteUser(c *gin.Context) {
	//on recupère l'id
	userID := c.Param("id")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"id": userID})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Can't find the user!!"})
	} else {
		if result.DeletedCount > 0 {
			fmt.Println("Count : " + string(result.DeletedCount))
			fmt.Println("USER WITH ID : " + userID + " HAS BEEN DELETED")
			//la suppression du fichier
			//ici je vais supprimer le fichier généré lors de l'insertion de l'utilisateur
			fileErr := os.Remove("./files/" + userID)
			if fileErr != nil {
				fmt.Println("Fichier non existant ou deja supprime")
			} else {
				fmt.Println("Le fichier " + userID + " a ete supprime avec succes")
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "L'utilisateur a ete supprime"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Can't find the user!!"})
		}

	}

}

//la fonction pour modifier un utilisateur
func EditUserByID(c *gin.Context) {
	//on recupère l'id
	userID := c.Param("id")
	fmt.Println("ID TO EDIT : " + userID)
	//je récupére le payload qui contient les données en json
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

	var newData interface{}
	//on recupère les données, la conversion en STRUCT
	err = json.Unmarshal(jsonData, &newData)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(newData)
	//la mise-à-jour des données dans la base de données
	update := bson.M{
		"$set": newData,
	}
	results, updateErr := collection.UpdateOne(context.TODO(), bson.M{"id": userID}, update)
	if updateErr != nil {
		fmt.Println("Can't edit the user")
		fmt.Println(updateErr.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Can't update the user!!"})
	} else {
		if results.ModifiedCount > 0 {
			fmt.Println("L'utilisateur a ete bien modifie")
			var user models.User
			err := collection.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&user)
			//si l'erreur est != nil ===> utilisateur non existant
			if err != nil {
				//vérification si le data a été changé ou pas (si oui, on va mettre à jour le fichier de l'utilisateur)
				content, err := ioutil.ReadFile("./files/" + userID)
				if err != nil {
					log.Fatal(err)
				}

				// conevrsion en string
				data := string(content)
				//si le nv data est diff de celui stocké dans le fichier on fait la m-à-j
				if user.Data != data {
					//écraser le l'ancien fichier et la création d'un nv contenant le nv data
					err := ioutil.WriteFile("./files/"+user.ID, []byte(user.Data), 0755)
					if err != nil {
						fmt.Printf("Unable to write file: %v", err)
					} else {
						fmt.Printf("Fichier " + user.ID + " a ete modifie avec succés.")
					}
				}

			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User infos updated"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Can't update the user!!"})
		}
	}

}
