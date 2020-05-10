# go-aws-zipper
A go project to create zip files from AWS S3

# Request example (JSON)
``
{
	"bucket": "bucket-name",
	"zip_name": "myZip",
	"items": [
		"somefile.png",
		"folder/anotherfile.png"
	]
}
``