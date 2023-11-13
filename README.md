# Gorage
A simple to use local storage system, which uses json files to store data and provides an easy to use go module to interacte with the data
<center>
<img src="https://i.imgur.com/8HDAwXt.png" alt="drawing" width="300"/>
</center>

## Features
- [X] Basic eval for where condition 
- [X] Insert statement
- [X] Create statement
- [X] Select statement
- [X] Update statement
- [X] Delete statement
- [ ] Eval Safety
- [X] Concurrency 
- [X] Advanced eval:  (, ), NAND, NOR

## Concurrency
GorageTable operations that are designed to be used concurrent:
- Insert
- Update 
- Delete

Every other operation can be used concurrent, but it is not guaranteed, that it will work.


## Create Storage and Tables
### CreateNewGorage
> `CreateNewGorage("./test", false, true)`
1. Paramter is the path and the file, which you want to create
2. Paramter is a boolean, if duplicate rows shall be allowed
3. Paramter is a boolean, if you want to see the log

### OpenGorage
> `g := OpenGorage("./test")`

Open Gorage by path
### CreateTable
> `g := OpenGorage("./test.json")`
> 
> `table := g.CreateTable("Example")`
1. Open Gorage

### AddColumn

```go
g := OpenGorage("./test.json")
table := g.CreateTable("User")
if table != nil {
	table.AddColumn("FirstName", STRING).
		AddColumn("LastName", STRING).
		AddColumn("Age", INT).
		AddColumn("IQ", INT)
}
```

## Data Operations
### FromTable
The `FromTable` function takes a string, which is the name of the table and returns a pointer to that table
### Delete
The `Delete` function takes no parameters and should only be called after the `.Where()`call, except you want to delete every row.

The delete function directly write to the real table in the memory.

Use `.Save()` to save it to the file for permanent change.

### Insert
The `Insert` function takes an array of `[]interface{}`. This list has to be the same length as the columns.

If you want to leave a cell blank just use `nil`.

Use `.Save()` to save it to the file for permanent change.

### Update
The `Update` function take a `map[string]interface{}`, where string is the column and interface is the new value.

The new data for the column needs to match the datatype, which the column can represent.

### Select
The `Select` function takes an array of strings, which represent the column names and returns a table, which only contains these columns.

This table is NOT persistent
### Where
The `Where` function take a string and can be used on a table to apply a filter. To compare data from the rows you can use :(Column).

This table is NOT persistent
#### Example
Let's say we have this table:


Name | Age | Country
--|---|---|
William | 20 | England
William | 22 | USA

If we now want to apply a filter to retrieve the rows where the name is 'William' we can do:
> `":Name == 'William' `

If we now want to apply a filter to retrieve the rows where the name is 'William' and the country is england we can do:
>`":Name == 'William' && :Country == 'England'`

See Eval Operations for syntax and operators

### Examples

Let's say we have this table:

Name | Age | Country
--|---|---|
William | 20 | England
William | 22 | USA

#### Select
```go
g := OpenGorage("./test.json")
userTable := g.FromTable("User").Where(":Name == 'William' && :Country == 'USA' ").Select([]string{"Name", "Age"})
```

### Update
```go
g := OpenGorage("./test.json")
g.FromTable("User").Where(":Name == 'William' && :Age == 20").Update(map[string]interface{}{
	"Name": "Tom"
})
g.Save()

```

#### Delete
```go
g := OpenGorage("./test.json")
g.FromTable("User").Where(":Name == 'William' && :Age == 20").Delete()
g.Save()
```

#### Insert
```go
g := OpenGorage("./test.json")
userTable := g.FromTable("User")
userTable.Insert([]interface{}{"Thomas", 33, nil})
userTable.Insert([]interface{}{"Carlos", 55, "USA"})
userTable.Insert([]interface{}{"Anna", nil, "USA"})
g.Save()
```

## Eval Operations
Eval currently supports && (AND), || (OR), !& (NAND), !| (NOR), == (EQUAL), != (NOT EQUAL), <, <=, >, >= and Braces ( ). 

**Important**: Spaces are important!

> Example 1: `"a == 5 || ( name == 'William' || name == 'James' )"`
> 
> Here is checked, if a is 5 or if the name is William or James


> Example 1: `"a == 5 !| b == 10"`
>
> Here is checked, if not a is 5 or b is 10


