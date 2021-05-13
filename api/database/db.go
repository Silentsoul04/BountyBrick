package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/splinter0/api/models"
)

/*
	Here are all the functions used to interact with the
	database.
*/

func DBInstance() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	URI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("bountybrick").Collection(collectionName)
}

/* USER MANAGEMENT */

func AddUser(user models.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	_, err := OpenCollection(Client, "users").InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func FindUser(username string) models.User {
	var u models.User
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	users := OpenCollection(Client, "users")
	users.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	defer cancel()
	return u
}

func AddUserToken(username, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	users := OpenCollection(Client, "users")
	users.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.D{
			{"$set", bson.D{{"last", time.Now()}, {"token", token}}},
		},
	)
	defer cancel()
}

func GetUserToken(username string) string {
	user := FindUser(username)
	return user.Token
}

/* REPO MANAGEMENT */

func AddRepo(repo models.Repo) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	_, err := OpenCollection(Client, "repos").InsertOne(ctx, repo)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Checks if repo is already in DB by link
func RepoExistsLink(link string) (bool, primitive.ObjectID) {
	var repo models.Repo
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	err := OpenCollection(Client, "repos").FindOne(ctx, bson.M{"link": link}).Decode(&repo)
	defer cancel()
	if err != nil {
		return false, repo.ID
	}
	return true, repo.ID
}

// Checks if repo is already in DB by ID
func RepoExistsId(id primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	result := OpenCollection(Client, "repos").FindOne(ctx, bson.M{"_id": id})
	defer cancel()
	return result.Err() == nil
}

func GetRepo(id primitive.ObjectID) models.Repo {
	var repo models.Repo
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	OpenCollection(Client, "repos").FindOne(ctx, bson.M{"_id": id}).Decode(&repo)
	defer cancel()
	return repo
}

func GetRepoByName(name string) models.Repo {
	var repo models.Repo
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	OpenCollection(Client, "repos").FindOne(ctx, bson.M{"name": name}).Decode(&repo)
	defer cancel()
	return repo
}

func GetAllRepos() []models.Repo {
	var repos []models.Repo
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	p := OpenCollection(Client, "repos")
	cursor, err := p.Find(ctx, bson.D{})
	if err == nil {
		for cursor.Next(ctx) {
			var r models.Repo
			cursor.Decode(&r)
			repos = append(repos, r)
		}
	}
	defer cancel()
	return repos
}

func EditRepo(id primitive.ObjectID, change bson.M) {
	repos := OpenCollection(Client, "repos")
	updated, err := repos.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{
			"$set": change,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(updated.MatchedCount, updated.ModifiedCount)
}

func SetForked(id primitive.ObjectID, status bool) {
	EditRepo(id, bson.M{"forked": status})
}

func SetBrick(id primitive.ObjectID, brickID string) {
	EditRepo(id, bson.M{"brick": brickID})

}

func SetCommit(id primitive.ObjectID, commit string) {
	EditRepo(id, bson.M{"commit": commit})

}

func SetStatus(id primitive.ObjectID, status string) {
	EditRepo(id, bson.M{"scan_status": status})

}

func SetResult(id primitive.ObjectID, result string) {
	EditRepo(id, bson.M{"scan_result": result})

}

func AddVulns(id primitive.ObjectID, vulns []models.Vuln) {
	repos := OpenCollection(Client, "repos")
	repos.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"vulns": bson.M{"$each": vulns}}},
	)
}

/* PROGRAMS MANAGEMENT */

func AddProgram(program models.Program) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	_, err := OpenCollection(Client, "programs").InsertOne(ctx, program)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetAllPrograms() []models.Program {
	var programs []models.Program
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	p := OpenCollection(Client, "programs")
	cursor, err := p.Find(ctx, bson.D{})
	if err == nil {
		for cursor.Next(ctx) {
			var prog models.Program
			cursor.Decode(&prog)
			programs = append(programs, prog)
		}
	}
	defer cancel()
	return programs
}

func GetProgramByName(name string) models.Program {
	var prog models.Program
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	OpenCollection(Client, "programs").FindOne(ctx, bson.M{"name": name}).Decode(&prog)
	defer cancel()
	return prog
}

func GetProgramById(id primitive.ObjectID) models.Program {
	var prog models.Program
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	OpenCollection(Client, "programs").FindOne(ctx, bson.M{"_id": id}).Decode(&prog)
	defer cancel()
	return prog
}

func UpdateProgramRepos(id primitive.ObjectID, repos []primitive.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	programs := OpenCollection(Client, "programs")
	programs.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"repos", repos}}},
		},
	)
	defer cancel()
}
