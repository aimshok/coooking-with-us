
async function goToReg(){
    window.location.href = "/registerPage";
}
async function handleOAuthCallback() {
    try {
        const response = await fetch('/auth/callback', {
            method: 'GET',
            credentials: 'include', // Include cookies for session management
        });

        if (!response.ok) {
            console.error(`Failed to fetch: ${response.statusText}`);
            return;
        }

        const data = await response.json();
        if (data.redirect) {
            // Perform the redirect if the backend provides the URL
            window.location.href = data.redirect;
        } else {
            console.error('No redirect URL provided by the backend.');
        }
    } catch (error) {
        console.error('Error handling OAuth callback:', error);
    }
}

// Call this function when the user is redirected back to your frontend
handleOAuthCallback();

async function login() {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    // Validation
    if (!email || !password) {
        alert("Email and password are required.");
        return;
    }

    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        const result = await response.json();

        if (result.status === "success") {
            alert("You logged in successfully");
            window.location.href = "/mainPage"; // Redirect to the main page
        } else {
            alert(result.message);
        }
    } catch (error) {
        alert("An error occurred: " + error.message);
    }
}