//https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/go/example_code

package clients
// import (
//     "github.com/aws/aws-sdk-go/aws"
//     "github.com/aws/aws-sdk-go/aws/session"
//     "github.com/aws/aws-sdk-go/service/s3"
//     "log"
// 	"strings"
// 	"github.com/oxygen-org/client/utils"
// )

// var locations = utils.StrSlice{"cn-north-1", "cn-northwest-1"}

// func GetS3SVC() *s3.S3{
// 	sess, err := session.NewSession(&aws.Config{
//         Region: aws.String("us-west-2")},
// 	)
// 	// Create S3 service client
// 	svc := s3.New(sess)
// 	return svc
// }

// func S3CreateBucket(bucket ,location string) string{
// 	if !locations.Contains(location){
// 		return "location must in cn-north-1,cn-northwest-1"
// 	}
// 	return ""

// }

// func S3ListBuckets() {

// }

// func S3GetAllObjects() {

// }

// func S3UploadFile() {

// }

// func S3DownloadFile() {

// }


func S3DeleteFile(){
	
}