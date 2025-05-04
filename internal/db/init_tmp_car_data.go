package carDataDB

func InitTmpCarData() {
	InsertCarData("1.png", 0.1, -30)
	InsertCarData("2.png", 0.2, -20)
	InsertCarData("3.png", 0.3, -10)
	InsertCarData("4.png", 0.4, 0)
	InsertPredictedCarData("5.png", 0.5, 10)
	InsertPredictedCarData("6.png", 0.6, 20)
	InsertPredictedCarData("7.png", 0.7, 30)
	InsertPredictedCarData("8.png", 0.8, 40)
	InsertPredictedCarData("9.png", 0.9, 50)
}
