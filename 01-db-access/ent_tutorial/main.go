
package main

import (
	"context"
	"log"
	"spike/01-db-access/ent_tutorial/ent"
	"spike/01-db-access/ent_tutorial/ent/car"
	"spike/01-db-access/ent_tutorial/ent/group"
	"spike/01-db-access/ent_tutorial/ent/user"

	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	client, err := ent.Open("postgres", "host=localhost port=5432 user=spike dbname=spike password=spike sslmode=disable")
	if err != nil {
        log.Fatalf("failed opening connection to postgres: %v", err)
    }
	defer client.Close()
	if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }

	// CreateUser(context.Background(), client)
	// QueryUser(context.Background(), client)
	// CreateCars(context.Background(), client)

	// QueryCars(context.Background(), a8m[0])
    // QueryCars(context.Background(), a8m[3])

	// QueryCarUsers(context.Background(), a8m[3])

	// if err := CreateGraph(context.Background(), client); err != nil {
    //     log.Fatalf("failed creating graph: %v", err)
    // }

	// QueryGitHub(context.Background(), client)
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
    u, err := client.User.
        Create().
        SetAge(30).
        SetName("a8ma").
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating user: %w", err)
    }
    log.Println("user was created: ", u)
    return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
    u, err := client.User.
        Query().
        Where(user.Name("a8ma")).
        // `Only` fails if no user found,
        // or more than 1 user returned.
        Only(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed querying user: %w", err)
    }
    log.Println("user returned: ", u)
    return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
    // Create a new car with model "Tesla".
    tesla, err := client.Car.
        Create().
        SetModel("Tesla").
        SetRegisteredAt(time.Now()).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating car: %w", err)
    }
    log.Println("car was created: ", tesla)

    // Create a new car with model "Ford".
    ford, err := client.Car.
        Create().
        SetModel("Ford").
        SetRegisteredAt(time.Now()).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating car: %w", err)
    }
    log.Println("car was created: ", ford)

    // Create a new user, and add it the 2 cars.
    a8m, err := client.User.
        Create().
        SetAge(30).
        SetName("a8m").
        AddCars(tesla, ford).
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating user: %w", err)
    }
    log.Println("user was created: ", a8m)
    return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
    cars, err := a8m.QueryCars().All(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    log.Println("returned cars:", cars)

    // What about filtering specific cars.
    ford, err := a8m.QueryCars().
        Where(car.Model("Ford")).
        Only(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    log.Println(ford)
    return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
    cars, err := a8m.QueryCars().All(ctx)
    if err != nil {
        return fmt.Errorf("failed querying user cars: %w", err)
    }
    // Query the inverse edge.
    for _, c := range cars {
        owner, err := c.QueryOwner().Only(ctx)
        if err != nil {
            return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
        }
        log.Printf("car %q owner: %q\n", c.Model, owner.Name)
    }
    return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	ariel, err := client.User.
        Create().
        SetAge(30).
        SetName("Ariel").
        Save(ctx)
    if err != nil {
        return err
    }

	neta, err := client.User.
		Create().
		SetAge(28).
		SetName("Neta").
		Save(ctx)
	if err != nil {
		return err
	}

	err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		SetOwner(ariel).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		SetOwner(ariel).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		SetOwner(neta).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(ariel).
		Exec(ctx)
	if err != nil {
        return err
    }

	err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(ariel, neta).
		Exec(ctx)
	if err != nil {
        return err
    }

    log.Println("The graph was created successfully")
    return nil
}

func QueryGitHub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.Name("GitHub")).
		QueryUsers().
		QueryCars().
		All(ctx)

	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
    }
    log.Println("cars returned:", cars)
    return nil
}