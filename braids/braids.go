package braids

import (
	"math"
	"math/rand"
)

/*
1	2	3	4
|	|	|	|
1	2	3	4 	[1,2,3,4]
\	/	|	|
|	|	|	|
1	2	3	4	[2,1,3,4] {1 -> 2, 2->1, 3->3, 4->4} (1 над 2 или 2 над 1 - хранит соответствующий номер записи в nodes)

*/

//Braid
/*
Мы представляем косу как "список пересечений" (Artin presentation).

Вместе с этим, в косе так же хранится информация о том, куда в итоге приходит отображение.

Вместо map используется list, т.к. это одно и то же в нашей ситуации,
т.к. у нас отображение из чисел от 1 до n без пробелов

*/
type Braid struct {
	nodes    []int // Перечень "Транзакций"
	width    int   //Количество "струн"(нитей) косы
	finalMap []int //Финальное отображение косы. Т.е.
}

func (left *Braid) Width() int {
	return left.width
}

// GenEmptyBraid
//Возвращает пустую косу заданного размера
func GenEmptyBraid(width int) Braid {
	trueMap := make([]int, width)
	for i := 0; i < width; i++ {
		trueMap[i] = i + 1 //1 -> 1, 2 -> 2, etc.
	}
	return Braid{
		make([]int, 0),
		width,
		trueMap,
	}
}

//Генерирует косу, заполненную случайными числами в размере nodesCount
func GenRandomBraid(width, nodesCount int) Braid {
	braid := GenEmptyBraid(width)
	for i := 0; i < nodesCount; i++ {
		number := -width + (rand.Int() % ((width - 1) * 2)) + 1
		braid.AddNode(number)
	}
	return braid
}

func GenHandmadeBraid(width int, nodes ...int) Braid {
	braid := GenEmptyBraid(width)
	braid.AddNode(nodes...)
	return braid
}

func (left *Braid) AddNode(nodes ...int) {
	for _, node := range nodes {
		abs := int(math.Abs(float64(node)))
		//Узел может быть только +-(1..n-1). 0-узел считается за eps и по-умолчанию исключается.
		if abs > 0 && abs < left.width {
			left.nodes = append(left.nodes, node)
			//Ну и выполняем перемещение, соответственно.
			left.finalMap[abs], left.finalMap[abs-1] = left.finalMap[abs-1], left.finalMap[abs]
		}

	}
}

//Возвращает "Обратную" версию косы. x.Reversed() = x^-1
func (left Braid) Reversed() Braid {
	res := Braid{
		make([]int, len(left.nodes)),
		left.width,
		make([]int, left.width),
	}
	copy(res.finalMap, left.finalMap)
	for i := len(left.nodes) - 1; i >= 0; i-- {
		res.nodes = append(res.nodes, -left.nodes[i])
	}
	return res
}

func (left Braid) Copy() Braid {
	res := Braid{
		make([]int, len(left.nodes)),
		left.width,
		make([]int, left.width),
	}
	copy(res.nodes, left.nodes)
	copy(res.finalMap, left.finalMap)

	return res
}

/*Выполняется "причёсывание" косы, включающее в себя избавление от пар типа x x^-1
 */
func (left *Braid) Brush() {
	res := GenEmptyBraid(left.width)
	brushed := make([]int, 0)
	for _, v := range left.nodes {
		if v != 0 {
			if len(brushed) == 0 {
				brushed = append(brushed, v)
			} else if brushed[len(brushed)-1] == v*-1 {
				brushed = brushed[:len(brushed)-1]
				continue
			} else {
				brushed = append(brushed, v)
			}
		}
	}
	//Восстанавливаю корректность итогового отображения
	for _, v := range brushed {
		res.AddNode(v)
	}
	left.nodes = brushed
	left.finalMap = res.finalMap
}

// Mult Представляет собой реализацию умножения на косу справа
func Mult(left Braid, elems ...Braid) Braid {
	res := left.Copy()
	for _, v := range elems {
		res = mult(res, v)
	}
	return res
}

// Mult Представляет собой реализацию умножения на косу справа
func mult(left, right Braid) Braid {
	res := left.Copy()
	for i, to := range right.finalMap {
		// Это похоже на композицию перестановок, но снизу-вверх. Оно корректно работает, честно.
		res.finalMap[i] = left.finalMap[to-1]
	}
	res.nodes = append(res.nodes, right.nodes...) //Добавляем все переходы из правого в левый
	res.Brush()
	return res
}

// x^n = {x * x} n times
func Pow(left Braid, n int) Braid {
	res := left.Copy()
	for i := 0; i < n; i++ {
		res = Mult(res, left)
	}
	return res
}

//Вспомогательная функция для сравнения двух множеств узлов
func nodeEquals(left, right []int) bool {
	for i, v := range left {
		if right[i] != v {
			return false
		}
	}
	return true
}
func (left Braid) Equals(right Braid) bool {
	return left.width == right.width && nodeEquals(right.nodes, left.nodes)
}
