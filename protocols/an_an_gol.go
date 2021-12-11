package protocols

import (
	"fmt"
	"go_braids/braids"
	"os"
)

const Width = 5
const Nodes = 8

//Та самая f()
func aliceAggregate(params ...braids.Braid) braids.Braid {
	res := braids.GenEmptyBraid(params[0].Width())
	res = braids.Mult(res, params...)
	return res
}

//Та самая t()
func bobAggregate(params ...braids.Braid) braids.Braid {
	res := braids.GenEmptyBraid(params[0].Width())
	for i, b := range params {
		if i%2 == 0 {
			res = braids.Mult(res, b)
		} else {
			res = braids.Mult(res, b.Reversed())
		}
	}
	return res
}

//Вспомогательная функция для генерации n кос
func getRandomBraids(amount int) []braids.Braid {
	res := make([]braids.Braid, amount)
	for i := 0; i < amount; i++ {
		res[i] = braids.GenRandomBraid(Width, Nodes)
	}
	return res
}

//Эпилог, вынесенный из кода, чтобы не захламлять всё подряд
func epilog(aliceKey, bobKey braids.Braid) {
	if aliceKey.Equals(bobKey) {
		fmt.Println("Обмен сообщениями прошёл успешно! Алиса и Боб установили связь и сформировали общий секрет!")
	} else {
		fmt.Println("Ошибка! Общий секрет не установлен!")
	}
	fmt.Println("В главных ролях:")
	fmt.Println("Расшифровка Алисы: ", aliceKey)
	fmt.Println("Расшифровка Боба: ", bobKey)
	fmt.Scanln()
}

//Вспомогательная функция, которая делает y^-1 * x * y для каждого члена последовательности
func genSeqForTilda(multer braids.Braid, params ...braids.Braid) []braids.Braid {
	res := make([]braids.Braid, len(params))
	for i, v := range params {
		res[i] = braids.Mult(multer.Reversed(), v, multer)
	}
	return res
}
func AnshelAnshelGoldfel() {
	var (
		aliceNumber, bobNumber int
	)
	fmt.Print("Реализация протокола Аншеля-Аншеля-Гольдфельда.\n\n")
	aliceFunc := aliceAggregate
	bobFunc := bobAggregate
	fmt.Print("Введите кол-во кос Алисы: ")
	fmt.Fscan(os.Stdin, &aliceNumber)
	fmt.Scanln()
	fmt.Print("Введите кол-во кос Боба: ")
	fmt.Fscan(os.Stdin, &bobNumber)
	fmt.Scanln()

	alicePubSet := getRandomBraids(aliceNumber)
	bobPubSet := getRandomBraids(bobNumber)
	//Алиса и Боб передали друг другу свои множества кос
	x := aliceFunc(bobPubSet...)
	y := bobFunc(alicePubSet...)
	//Алиса и Боб сформировали публичные ключи, основываясь на полученной публичной информации друг от друга.
	//Ключи получены при помощи функции-секрета
	xTilda := aliceFunc(genSeqForTilda(y, bobPubSet...)...)
	yTilda := bobFunc(genSeqForTilda(x, alicePubSet...)...)
	aliceKey := braids.Mult(x.Reversed(), xTilda)
	bobKey := braids.Mult(yTilda.Reversed(), y)
	epilog(aliceKey, bobKey)

}
