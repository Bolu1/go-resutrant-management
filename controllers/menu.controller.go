package controllers

import(
	"resturant-management/database"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"resturant-management/models"
	"context"
	"log"
	"time"
	"net/http"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while lisiting the menu items"})
		}
		var allMenu []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenu)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		fmenuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while fetching the menu item"})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = cotext.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		validationErr := validate.Struct(menu)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validatonErr.Error()})
		}
		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr = menuCollection.InsertOne(ctx, menu)
		if insertErr != {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()

	}
}
func inTimeSpan(start, end, check time.Time) bool{
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.menu

		if err := c.BindJSON(&menu); err !=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		menuId := c.Params("menu_id")
		filter := bson.M{"menu_id":menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil{
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()){
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
			updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})

			if menu.Name != ""{
				updateObj = append(updateObj, bson.E{"name", menu.Name})
			}
			if menu.Category != ""{
				updateObj = append(updateObj, bson.E{"category", menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"update_at", menu.Update_at})

			upsert := true

			opt := options.UpdateOptions{
				Upsert : &upsert
			}

			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set", updateObj},
				},
				&opt,
			)

			if err != nil{
				msg := "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}

			defer cancel()
			c.JSON(http.StatusOK, result)
			
		}
	}
}

