document.getElementById('registrationForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const name = document.getElementById('name').value
    if (password !== confirmPassword) {
        document.getElementById('message').innerText = 'Passwords do not match.';
        return;
    }

    const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        }

    );

    const result = await response.json();
    document.getElementById('message').innerText = result.message;

    if (response.ok) {
        localStorage.setItem('email', email);
        localStorage.setItem('name', name);
        localStorage.setItem('password', password);// Сохраняем email для подтверждения
        window.location.href = '/verifyEmailPage';
    }
});