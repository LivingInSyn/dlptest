BuildDocker:
	docker build -t dlptest .

Run:
	docker run -it -e S3REGION='us-west-2' -e S3BUCKET='livinginsyn-dlptest' -e S3KEYID='${DLPTEST_AK}' -e S3SECRET='${DLPTEST_SK}' -p 8080:8080 dlptest

RunInteractive:
	docker run -it -e S3REGION='us-west-2' -e S3BUCKET='livinginsyn-dlptest' -e S3KEYID='${DLPTEST_AK}' -e S3SECRET='${DLPTEST_SK}' dlptest /bin/sh