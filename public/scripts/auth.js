document.addEventListener('DOMContentLoaded', function() {
    // Получаем формы
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    
    // Обработка формы входа
    if (loginForm) {
        loginForm.addEventListener('submit', async function(event) {
            event.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const errorDiv = document.getElementById('login-error');
            
            try {
                const response = await fetch('/api/users/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ username, password })
                });
                
                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || 'Ошибка входа');
                }
                
                const userData = await response.json();
                
                // Сохраняем данные пользователя в localStorage
                localStorage.setItem('user', JSON.stringify({
                    id: userData.id,
                    username: userData.username,
                    token: userData.token,
                    fullName: userData.full_name
                }));
                
                // Перенаправляем на страницу заказов
                window.location.href = '/order';
            } catch (error) {
                errorDiv.textContent = error.message;
                errorDiv.style.display = 'block';
            }
        });
    }
    
    // Обработка формы регистрации
    if (registerForm) {
        registerForm.addEventListener('submit', async function(event) {
            event.preventDefault();
            
            const username = document.getElementById('username').value;
            const email = document.getElementById('email').value;
            const fullName = document.getElementById('full_name').value;
            const password = document.getElementById('password').value;
            const confirmPassword = document.getElementById('confirm_password').value;
            const errorDiv = document.getElementById('register-error');
            
            // Проверяем совпадение паролей
            if (password !== confirmPassword) {
                errorDiv.textContent = 'Пароли не совпадают';
                errorDiv.style.display = 'block';
                return;
            }
            
            try {
                const response = await fetch('/api/users/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        username,
                        email,
                        full_name: fullName,
                        password
                    })
                });
                
                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.error || 'Ошибка регистрации');
                }
                
                const userData = await response.json();
                
                // Сохраняем данные пользователя в localStorage
                localStorage.setItem('user', JSON.stringify({
                    id: userData.id,
                    username: userData.username,
                    token: userData.token,
                    fullName: userData.full_name
                }));
                
                // Перенаправляем на страницу заказов
                window.location.href = '/order';
            } catch (error) {
                errorDiv.textContent = error.message;
                errorDiv.style.display = 'block';
            }
        });
    }
    
    // Проверка авторизации пользователя
    function checkAuth() {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        
        // Если пользователь авторизован и находится на странице авторизации, 
        // перенаправляем на страницу заказов
        if (user.token && (window.location.pathname === '/login' || window.location.pathname === '/register')) {
            window.location.href = '/order';
        }
    }
    
    // Проверяем авторизацию при загрузке страницы
    checkAuth();
});