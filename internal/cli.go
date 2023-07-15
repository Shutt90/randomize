package cli

// TODO: Rearchitect this cli stuff
// func storeCli() {
// 	fmt.Println("Enter website name: ")
// 	var entry storedPassword
// 	_, err := fmt.Scanln(&entry.WebsiteName)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	_, err = fmt.Scanln(&entry.Username)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Number of characters")
// 	var chars string
// 	n, err := fmt.Scanln(&chars)
// 	if err != nil || n > 3 {
// 		fmt.Println("Error: can't be higher than 255 and ", err)
// 		return
// 	}
// 	convertedNum, err := strconv.Atoi(chars)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	if convertedNum > 255 || convertedNum < 8 {
// 		fmt.Println("Error: can't be higher than 255 or lower than 8")
// 		return
// 	}

// 	entry.Password = helpers.Randomize(uint8(convertedNum))
// 	fmt.Println("your password is ", entry.Password)

// 	// entry.store(ctx)
// }

// func getCli() {
// 	fmt.Println("Enter website name: ")
// 	var websiteName string
// 	_, err := fmt.Scanln(&websiteName)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	ctx := context.Background()
// 	pw, err := getPassword(websiteName, ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(pw)
// }
