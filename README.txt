go mod init kp-app

   gin = https://github.com/gin-gonic/gin    


.vscode 
{
    "go.useLanguageServer": true,
    "[go]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true,
        },
        // Optional: Disable snippets, as they conflict with completion ranking.
        "editor.snippetSuggestions": "none",
    },
    "[go.mod]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true,
        },
    },
    "gopls": {
        // Add parameter placeholders when completing a function.
        "usePlaceholders": true,

        // If true, enable additional analyses with staticcheck.
        // Warning: This will significantly increase memory usage.
        "staticcheck": false,
    }
}

//upload file
// Get file
		file, _ := ctx.FormFile("image")

		// Create file
		path := "uploads/products/" + strconv.Itoa(int(p.ID)) // ID => 8, uploads/articles/8/image.png
		os.MkdirAll(path, 0755)                               // -> uploads/products/8

		// Upload file
		filename := path + file.Filename
		if err := ctx.SaveUploadedFile(file, filename); err != nil {
			log.Fatal(err.Error())
		}

		// Attach file to products
		p.Image = "http://localhost:8080/" + filename




// connect <i class="fa fa-database

func (p *Product) setProductImage(ctx *gin.Context, products *models.Product) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return nil
	}

	if products.Image == "" {
		//1. ตัด http://localhost:80800/uploads/prosucts/<ID>/image.png ให้เหลือ /uploads/prosucts/<ID>/image.png
		products.Image = strings.Replace(products.Image, os.Getenv("HOST"), "", 1)
		//2. แทนค่าพาธปัจจุบัน<WD>/uploads/prosucts/<ID>/image.png
		pwd, _ := os.Getwd()
		//3.remove <WD>/uploads/prosucts/<ID>/image.png
		os.Remove(pwd + products.Image)
	}

	path := "uploads/products/" + strconv.Itoa(int(products.ID))
	os.MkdirAll(path, 0755)

	filename := path + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	products.Image = os.Getenv("HOST") + "/" + filename

	p.DB.Save(products)

	return nil
}
      