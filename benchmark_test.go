package Gorage

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestBenchmark(t *testing.T) {
	if fileExists("./benchmark") {
		err := os.Remove("./benchmark")
		if err != nil {
			t.Fatalf("Error removing old test file")
			return
		}
	}
	gorage := CreateNewGorage("./benchmark", false, false)
	randomTable := gorage.CreateTable("Random")
	if randomTable == nil {
		t.Fatalf("Table was not created")
	}
	t1 := time.Now().Unix()

	for i := 0; i < 20; i++ {
		randomTable.AddColumn(fmt.Sprintf("Col%d", i), INT)
	}
	for i := 0; i < 10000; i++ {
		t.Log("Currently at ", i)
		randomTable.Insert([]interface{}{
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
			rand.Intn(1000-0) + 0,
		})
	}
	t.Log(fmt.Sprintf("Inserting %d took %ds", 10*200, time.Now().Unix()-t1))

	t1 = time.Now().Unix()
	gorage.Save()
	gorage = OpenGorage("./benchmark")
	t.Log(fmt.Sprintf("Saving and loading took %ds", time.Now().Unix()-t1))

	t1 = time.Now().Unix()
	for i := 0; i < 10; i++ {
		col1 := rand.Intn(20-0) + 0
		col2 := rand.Intn(20-0) + 0
		c1 := rand.Intn(1000-0) + 0
		c2 := rand.Intn(1000-0) + 0
		randomTable.Where(fmt.Sprintf(":Col%d == %d && :Col%d == %d", col1, c1, col2, c2))
	}
	t.Log(fmt.Sprintf("Getting 10 random rows took %ds", time.Now().Unix()-t1))

	t1 = time.Now().Unix()
	for i := 0; i < 10; i++ {
		c1 := rand.Intn(1000-0) + 0
		c2 := rand.Intn(200-0) + 0
		randomTable.Where(fmt.Sprintf(":Col1 == %d && :Col2 == %d", c1, c2)).Delete()
	}
	t.Log(fmt.Sprintf("Deleting 10 random rows took %ds", time.Now().Unix()-t1))

	t1 = time.Now().Unix()
	for i := 0; i < 10; i++ {
		c1 := rand.Intn(1000-0) + 0
		c2 := rand.Intn(200-0) + 0
		randomTable.Where(fmt.Sprintf(":Col1 == %d && :Col2 == %d", c1, c2)).Update(map[string]interface{}{
			"Col1": 0,
		})
	}
	t.Log(fmt.Sprintf("Updating 10 random rows took %ds", time.Now().Unix()-t1))
	gorage.Save()

}
