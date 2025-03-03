/* http_test.js */

async function fileClick(filename, uploadType = "normal") {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const posturl = "http://localhost:8080/upload";

    try {
        let result = await downloadAndPostFile(dlurl, posturl, uploadType);
        console.log("Result from server:", result);
    } catch (error) {
        console.error("Error during file download and posting:", error);
    }
}

async function downloadAndPostFile(fileUrl, postUrl, uploadType) {
    const filename = fileUrl.substring(fileUrl.lastIndexOf("/") + 1);
    const statusSpan = document.getElementById(`${filename}-status`);

    try {
        const blob = await downloadFileToBlob(fileUrl);

        if (uploadType === "secure") {
            await secureUpload(blob, filename, statusSpan);
        }

        console.log("File posted successfully");
    } catch (error) {
        if (statusSpan) {
            statusSpan.textContent = `POST not successful - ${error.message}`;
        }
        console.error("Error downloading or posting file:", error);
        throw error;
    }
}

async function downloadFileToBlob(fileUrl) {
    try {
        const response = await fetch(fileUrl);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.blob();
    } catch (error) {
        console.error("Error downloading file:", error);
        throw error;
    }
}

async function secureUpload(file, filename, statusSpan) {
    try {
        const response = await fetch("/generate-s3-url");
        if (!response.ok) throw new Error("Failed to retrieve S3 URL.");

        const data = await response.json();
        const uploadUrl = data.url;

        const uploadResponse = await fetch(uploadUrl, {
            method: "PUT",
            body: file,
            headers: { "Content-Type": file.type }
        });

        if (!uploadResponse.ok) throw new Error("Failed to upload file to S3.");

        if (statusSpan) {
            statusSpan.textContent = "File securely uploaded to S3!";
        }
    } catch (error) {
        if (statusSpan) {
            statusSpan.textContent = `Error during secure upload: ${error.message}`;
        }
        console.error("Error during secure upload:", error);
    }
}