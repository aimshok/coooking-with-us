document.getElementById('verificationForm').addEventListener('submit', async (event) => {
    event.preventDefault();

    const code = document.getElementById('verificationCode').value;
    const email = localStorage.getItem('email');
    const password = localStorage.getItem('password');
    const name = localStorage.getItem('name');
    if (code.length !== 6) {
        alert("The verification code must be 6 digits");
        return;
    }

    try {
        const response = await fetch('/verifyEmail', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, name, code, password})
        });

        const result = await response.json();

        if (response.ok) {
            // Код подтвержден успешно
            alert('Email verified successfully!');
            window.location.href = '/loginPage';
        } else {
            // Ошибка, например, неверный код
            document.getElementById('message').innerText = result.message || 'Verification failed';
        }
    } catch (error) {
        // Обработка сетевых ошибок
        document.getElementById('message').innerText = 'An error occurred. Please try again later.';
        console.error('Error verifying email:', error);
    }
});
