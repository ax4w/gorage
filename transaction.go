package Gorage

const (
	actionSelect       = 0
	actionWhere        = 1
	actionInsert       = 2
	actionDelete       = 3
	actionUpdate       = 4
	actionAddColumn    = 5
	actionDeleteColumn = 6
	actionDeleteTable  = 7
	actionExit         = 8
	actionFromTable    = 9
)

type data struct {
	action  int
	payload []interface{}
	c       chan *Table
}

func transactionManger(t *Table) {
	for {
		d := t.t.q.Head()
		for d != nil {
			switch d.action {
			case actionFromTable:
				d.c <- t
				break
			case actionExit:
				goto exit
			case actionSelect:
				col := d.payload[0].([]string)
				d.c <- t._select(col)
				break
			case actionWhere:
				where := d.payload[0].(string)
				d.c <- t.where(where)
				break
			case actionInsert:
				d.c <- t.insert(d.payload)
				break
			case actionDelete:
				ti := d.payload[0].(*Table)
				d.c <- ti.delete()
				break
			case actionUpdate:
				p := d.payload[0].(map[string]interface{})
				d.c <- t.Update(p)
				break
			case actionAddColumn:
				name := d.payload[0].(string)
				datatype := d.payload[1].(int)
				r := t.addColumn(name, datatype)
				d.c <- r
				break
			case actionDeleteColumn:
				name := d.payload[0].(string)
				d.c <- t.removeColumn(name)
				break
			case actionDeleteTable:
				break
			}
			t.t.q = t.t.q.Shift()
			d = t.t.q.Head()
		}
	}
exit:
	t.t.q = nil
}
