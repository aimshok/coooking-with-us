document.addEventListener("DOMContentLoaded", async () => {
    console.log('DOMContentLoaded');
    try {
        const regmet = await fetch("/getUserRegMet");
        if (!regmet.ok) {
            throw new Error("Failed to fetch user regmet");
        }

        const regData = await regmet.json();
        const reg_met = regData.reg_met
        console.log("API Response:", regData);

        if(regData.reg_met === "google"){
            document.getElementById("passText").style.display = "none";
            document.getElementById("passUpd").style.display = "none";
        }
    } catch (error) {
        console.error("Error fetching user avatar:", error);
    }
    try {
        // Fetch user avatar
        const avatarResponse = await fetch("/getUserAvatar");
        if (!avatarResponse.ok) {
            throw new Error("Failed to fetch user avatar");
        }

        const avatarData = await avatarResponse.json();
        const avatarUrl = avatarData.avatar;

        if (avatarUrl) {
            document.getElementById("userAvatar").src = avatarUrl;
        } else {
            console.error("Avatar not found");
        }
    } catch (error) {
        console.error("Error fetching user avatar:", error);
    }

    // Existing name and email fetching logic
    try {
        const response = await fetch("/getUserName");
        if (!response.ok) {
            throw new Error("Failed to fetch user name");
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
        const response = await fetch("/getUserEmail");
        if (!response.ok) {
            throw new Error("Failed to fetch user email");
        }

        const data = await response.json();
        const userEmail = data.email;

        if (userEmail) {
            document.getElementById("userEmail").textContent = userEmail;
        } else {
            console.error("Email not found");
        }
    } catch (error) {
        console.error("Error fetching user email:", error);
    }
});

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

document.getElementById('avatarInput').addEventListener('change', async () => {
    const avatarInput = document.getElementById('avatarInput');
    if (avatarInput.files.length === 0) return;

    console.log("Avatar selected:", avatarInput.files[0]);

    const formData = new FormData();
    formData.append('avatar', avatarInput.files[0]);

    try {
        const response = await fetch('/updateAvatar', {
            method: 'POST',
            body: formData,
        });

        const result = await response.json();
        const statusMessage = document.getElementById('statusMessage');

        // Очистить старые сообщения
        statusMessage.innerHTML = '';

        if (result.status === 'success') {
            const successMessage = document.createElement('div');
            successMessage.className = 'alert alert-success';
            successMessage.textContent = result.message;
            statusMessage.appendChild(successMessage);

            // Обновляем отображаемую аватарку
            document.getElementById('userAvatar').src = result.avatarUrl;
        } else {
            const failMessage = document.createElement('div');
            failMessage.className = 'alert alert-danger';
            failMessage.textContent = result.message;
            statusMessage.appendChild(failMessage);
        }
    } catch (error) {
        console.error('Error uploading avatar:', error);
    }
});

// Функция для триггера загрузки при клике на аватар
function triggerAvatarUpload() {
    document.getElementById('avatarInput').click();
}
document.getElementById('nameUpd').addEventListener('submit', async (event) => {
    event.preventDefault();
    console.log("nameUpd" + document.getElementById("nameUpd").value);
    const name = document.getElementById('newName').value;

    const response = await fetch('/updateName', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(name)
    })
    const result = await response.json();
    if (result.status === "success") {
        const successMessage = document.createElement('div');
        successMessage.className = 'alert alert-success';
        successMessage.textContent = result.message;
        document.getElementById('statusMessage').appendChild(successMessage);
    } else {
        const failMessage = document.createElement('div');
        failMessage.className = 'alert alert-danger';
        failMessage.textContent = result.message;
        document.getElementById('statusMessage').appendChild(failMessage);
    }
})


document.getElementById('passUpd').addEventListener('submit', async (event) => {
    event.preventDefault();
    console.log("PassUpd"+document.getElementById('newPassword').value);
    const oldPassword = document.getElementById('oldPassword').value;
    const newPassword = document.getElementById('newPassword').value;

    const response = await fetch('/updatePassword', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({newPassword, oldPassword})
    })
    const result = await response.json();
    if (result.status === "success") {
        const successMessage = document.createElement('div');
        successMessage.className = 'alert alert-success';
        successMessage.textContent = result.message;
        document.getElementById('statusMessage').appendChild(successMessage);
    } else {
        const failMessage = document.createElement('div');
        failMessage.className = 'alert alert-danger';
        failMessage.textContent = result.message;
        document.getElementById('statusMessage').appendChild(failMessage);
    }
})

