

async function logOut(){
    const response = await fetch(`/logout`, {
        method: 'POST',
        credentials: 'same-origin',
    });


    const result = await response.json();

    if (response.ok) {
        alert(result.message)
        window.location.href = "/loginPage";
    } else {
        alert("Error logging out: " + result.message);
    }
}
async function sendEmail() {
    const subject = document.getElementById("subject").value;
    const body = document.getElementById("body").value;
    const file = document.getElementById("file").files[0];

    if (!subject.trim() || !body.trim()) {
        alert("Subject and body are required.");
        return;
    }

    const emailData = {
        subject: subject,
        body: body,
        toEmail: "adil_groud@mail.ru",
    };

    const formData = new FormData();
    formData.append("emailData", JSON.stringify(emailData));
    if (file) {
        formData.append("file", file);
        formData.append("fileName", file.name);
    }

    const response = await fetch("/sendEmail", {
        method: "POST",
        body: formData
    });

    const result = await response.json();
    if (response.ok) {
        alert(result.message);
    } else {
        alert("An error occurred: " + result.message);
    }

}

document.addEventListener("DOMContentLoaded", async () => {
    try {
        const response = await fetch("/getUserStatus");
        if (!response.ok) {
            throw new Error("Failed to fetch user email");
        }

        const data = await response.json();
        const status = data.status;
        if (status === "admin") {
            document.getElementById("adminButton").style.display = "block";
        }
    } catch (error) {
        console.error("Error fetching user email:", error);
    }
    try {
        const response = await fetch("/getUserName");
        if (!response.ok) {
            throw new Error("Failed to fetch user email");
        }

        const data = await response.json();
        const userName = data.name;

        if (userName) {
            document.getElementById("userName").textContent = userName;
        } else {
            console.error("User not found");
        }
    } catch (error) {
        console.error("Error fetching user name:", error);
    }
    try {
        // Fetch user avatar
        const avatarResponse = await fetch("/getUserAvatar");
        if (!avatarResponse.ok) {
            throw new Error("Failed to fetch user avatar");
        }

        const avatarData = await avatarResponse.json();
        const avatarUrl = avatarData.avatar;

        // Check if the avatar URL points to a valid image file
        const imageResponse = await fetch(avatarUrl, { method: "HEAD" });

        if (imageResponse.ok) {
            // If the avatar exists, set it as the image source
            document.getElementById("avatar").src = avatarUrl;
        } else {
            // If the avatar does not exist, use the default avatar
            console.error("Avatar not found, using default avatar.");
            document.getElementById("avatar").src = "../avatars/003.png";
        }
    } catch (error) {
        console.error("Error fetching user avatar:", error);
        document.getElementById("avatar").src = "../avatars/001.png"; // Default avatar
    }


});