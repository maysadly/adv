// Состояние приложения
const state = {
    products: [],
    cart: [],
    currentPage: 1,
    perPage: 6,
    total: 0,
    filters: {
        name: '',
        minPrice: '',
        maxPrice: ''
    },
    userId: '',
    orders: []
};

// DOM-элементы
const DOM = {
    productsList: document.getElementById('products-list'),
    pagination: document.getElementById('pagination'),
    cartItems: document.getElementById('cart-items'),
    cartEmpty: document.getElementById('cart-empty'),
    cartCount: document.getElementById('cart-count'),
    cartTotal: document.getElementById('cart-total'),
    cartSection: document.getElementById('cart-section'),
    cartIcon: document.getElementById('cart-icon'),
    checkoutButton: document.getElementById('checkout-button'),
    orderForm: document.getElementById('order-form'),
    userIdInput: document.getElementById('user-id'),
    nameFilter: document.getElementById('name-filter'),
    minPriceFilter: document.getElementById('min-price'),
    maxPriceFilter: document.getElementById('max-price'),
    applyFiltersButton: document.getElementById('apply-filters'),
    resetFiltersButton: document.getElementById('reset-filters'),
    userIdFilter: document.getElementById('user-id-filter'),
    viewOrdersButton: document.getElementById('view-orders'),
    ordersList: document.getElementById('orders-list')
};

// Инициализация приложения
document.addEventListener('DOMContentLoaded', () => {
    // Загрузка продуктов
    fetchProducts();

    // Обработчики событий
    DOM.applyFiltersButton.addEventListener('click', applyFilters);
    DOM.resetFiltersButton.addEventListener('click', resetFilters);
    DOM.orderForm.addEventListener('submit', placeOrder);
    DOM.viewOrdersButton.addEventListener('click', fetchOrders);
    DOM.userIdInput.addEventListener('input', updateCheckoutButton);

    // Проверяем, сохранена ли корзина в localStorage
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
        state.cart = JSON.parse(savedCart);
        updateCart();
    }

    // Проверяем, сохранен ли ID пользователя в localStorage
    const savedUserId = localStorage.getItem('userId');
    if (savedUserId) {
        DOM.userIdInput.value = savedUserId;
        DOM.userIdFilter.value = savedUserId;
        state.userId = savedUserId;
        updateCheckoutButton();
        fetchOrders();
    }
});

// Получение продуктов с фильтрацией и пагинацией
async function fetchProducts() {
    try {
        const url = new URL('/api/products', window.location.origin);
        url.searchParams.append('page', state.currentPage);
        url.searchParams.append('per_page', state.perPage);

        if (state.filters.name) {
            url.searchParams.append('name', state.filters.name);
        }
        if (state.filters.minPrice) {
            url.searchParams.append('min_price', state.filters.minPrice);
        }
        if (state.filters.maxPrice) {
            url.searchParams.append('max_price', state.filters.maxPrice);
        }

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Ошибка при получении продуктов: ${response.statusText}`);
        }

        const data = await response.json();
        state.products = data.products;
        state.total = data.total;
        
        renderProducts();
        renderPagination();
    } catch (error) {
        console.error('Ошибка при загрузке продуктов:', error);
    }
}

// Отображение списка продуктов
function renderProducts() {
    DOM.productsList.innerHTML = '';

    if (state.products.length === 0) {
        DOM.productsList.innerHTML = '<div class="main__products-empty">Продукты не найдены</div>';
        return;
    }

    state.products.forEach(product => {
        const productElement = document.createElement('div');
        productElement.className = 'main__product-item';
        
        const isInCart = state.cart.some(item => item.id === product.ID);
        const isOutOfStock = product.Stock <= 0;
        
        productElement.innerHTML = `
            <div class="main__product-item-name">${product.Name}</div>
            <div class="main__product-item-price">${product.Price.toFixed(2)} ₸</div>
            <div class="main__product-item-stock">В наличии: ${product.Stock}</div>
            <button class="main__product-item-button" 
                    data-id="${product.ID}" 
                    data-name="${product.Name}" 
                    data-price="${product.Price}" 
                    data-stock="${product.Stock}"
                    ${isInCart || isOutOfStock ? 'disabled' : ''}
            >
                ${isInCart ? 'В корзине' : isOutOfStock ? 'Нет в наличии' : 'Добавить в корзину'}
            </button>
        `;
        
        const addButton = productElement.querySelector('button');
        if (!isInCart && !isOutOfStock) {
            addButton.addEventListener('click', addToCart);
        }
        
        DOM.productsList.appendChild(productElement);
    });
}

// Отображение пагинации
function renderPagination() {
    DOM.pagination.innerHTML = '';
    
    const totalPages = Math.ceil(state.total / state.perPage);
    if (totalPages <= 1) return;
    
    // Кнопка "Предыдущая"
    if (state.currentPage > 1) {
        const prevButton = document.createElement('button');
        prevButton.textContent = 'Назад';
        prevButton.addEventListener('click', () => {
            state.currentPage--;
            fetchProducts();
        });
        DOM.pagination.appendChild(prevButton);
    }
    
    // Кнопки с номерами страниц
    for (let i = 1; i <= totalPages; i++) {
        const pageButton = document.createElement('button');
        pageButton.textContent = i;
        pageButton.disabled = i === state.currentPage;
        pageButton.addEventListener('click', () => {
            state.currentPage = i;
            fetchProducts();
        });
        DOM.pagination.appendChild(pageButton);
    }
    
    // Кнопка "Следующая"
    if (state.currentPage < totalPages) {
        const nextButton = document.createElement('button');
        nextButton.textContent = 'Вперед';
        nextButton.addEventListener('click', () => {
            state.currentPage++;
            fetchProducts();
        });
        DOM.pagination.appendChild(nextButton);
    }
}

// Добавление товара в корзину
function addToCart(event) {
    const button = event.target;
    const id = button.dataset.id;
    const name = button.dataset.name;
    const price = parseFloat(button.dataset.price);
    const stock = parseInt(button.dataset.stock);
    
    // Проверяем, есть ли уже этот товар в корзине
    const existingItem = state.cart.find(item => item.id === id);
    
    if (existingItem) {
        existingItem.quantity += 1;
    } else {
        state.cart.push({
            id,
            name,
            price,
            stock,
            quantity: 1
        });
    }
    
    // Обновляем UI
    button.disabled = true;
    button.textContent = 'В корзине';
    
    // Сохраняем корзину в localStorage
    localStorage.setItem('cart', JSON.stringify(state.cart));
    
    // Обновляем отображение корзины
    updateCart();
}

// Обновление отображения корзины
function updateCart() {
    const totalItems = state.cart.reduce((sum, item) => sum + item.quantity, 0);
    DOM.cartCount.textContent = totalItems;
    
    // Если корзина пуста, показываем сообщение
    if (totalItems === 0) {
        DOM.cartEmpty.style.display = 'block';
        DOM.cartItems.innerHTML = '';
        DOM.cartTotal.textContent = '0 ₸';
        DOM.checkoutButton.disabled = true;
        return;
    }
    
    // Скрываем сообщение о пустой корзине
    DOM.cartEmpty.style.display = 'none';
    
    // Очищаем содержимое корзины перед обновлением
    DOM.cartItems.innerHTML = '';
    
    // Общая сумма
    let totalPrice = 0;
    
    // Добавляем элементы корзины
    state.cart.forEach(item => {
        const itemTotal = item.price * item.quantity;
        totalPrice += itemTotal;
        
        const cartItem = document.createElement('div');
        cartItem.className = 'main__cart-item';
        cartItem.innerHTML = `
            <div class="main__cart-item-details">
                <div class="main__cart-item-name">${item.name}</div>
                <div class="main__cart-item-price">${item.price.toFixed(2)} ₸ x ${item.quantity}</div>
            </div>
            <div class="main__cart-item-quantity">
                <button class="decrease" data-id="${item.id}" ${item.quantity <= 1 ? 'disabled' : ''}>-</button>
                <span>${item.quantity}</span>
                <button class="increase" data-id="${item.id}" ${item.quantity >= item.stock ? 'disabled' : ''}>+</button>
            </div>
            <button class="main__cart-item-remove" data-id="${item.id}">×</button>
        `;
        
        // Добавляем обработчики событий для кнопок
        const decreaseButton = cartItem.querySelector('.decrease');
        const increaseButton = cartItem.querySelector('.increase');
        const removeButton = cartItem.querySelector('.main__cart-item-remove');
        
        decreaseButton.addEventListener('click', () => decreaseQuantity(item.id));
        increaseButton.addEventListener('click', () => increaseQuantity(item.id));
        removeButton.addEventListener('click', () => removeFromCart(item.id));
        
        DOM.cartItems.appendChild(cartItem);
    });
    
    // Обновляем общую сумму
    DOM.cartTotal.textContent = `${totalPrice.toFixed(2)} ₸`;
    
    // Активируем кнопку оформления заказа, если указан ID пользователя
    updateCheckoutButton();
}

// Уменьшение количества товара в корзине
function decreaseQuantity(id) {
    const item = state.cart.find(item => item.id === id);
    if (!item) return;
    
    item.quantity -= 1;
    
    if (item.quantity <= 0) {
        removeFromCart(id);
    } else {
        localStorage.setItem('cart', JSON.stringify(state.cart));
        updateCart();
    }
}

// Увеличение количества товара в корзине
function increaseQuantity(id) {
    const item = state.cart.find(item => item.id === id);
    if (!item || item.quantity >= item.stock) return;
    
    item.quantity += 1;
    localStorage.setItem('cart', JSON.stringify(state.cart));
    updateCart();
}

// Удаление товара из корзины
function removeFromCart(id) {
    state.cart = state.cart.filter(item => item.id !== id);
    localStorage.setItem('cart', JSON.stringify(state.cart));
    
    // Обновляем кнопку добавления товара
    const productButton = document.querySelector(`button[data-id="${id}"]`);
    if (productButton) {
        productButton.disabled = false;
        productButton.textContent = 'Добавить в корзину';
    }
    
    updateCart();
}

// Применение фильтров
function applyFilters() {
    state.filters.name = DOM.nameFilter.value;
    state.filters.minPrice = DOM.minPriceFilter.value;
    state.filters.maxPrice = DOM.maxPriceFilter.value;
    state.currentPage = 1;
    fetchProducts();
}

// Сброс фильтров
function resetFilters() {
    DOM.nameFilter.value = '';
    DOM.minPriceFilter.value = '';
    DOM.maxPriceFilter.value = '';
    state.filters = { name: '', minPrice: '', maxPrice: '' };
    state.currentPage = 1;
    fetchProducts();
}

// Обновление кнопки оформления заказа
function updateCheckoutButton() {
    const userId = DOM.userIdInput.value.trim();
    state.userId = userId;
    
    if (userId && state.cart.length > 0) {
        DOM.checkoutButton.disabled = false;
    } else {
        DOM.checkoutButton.disabled = true;
    }
    
    // Сохраняем ID пользователя
    if (userId) {
        localStorage.setItem('userId', userId);
    }
}

// Оформление заказа
async function placeOrder(event) {
    event.preventDefault();
    
    const userId = DOM.userIdInput.value.trim();
    if (!userId || state.cart.length === 0) return;
    
    try {
        const orderData = {
            user_id: userId,
            items: state.cart.map(item => ({
                product_id: item.id,
                quantity: item.quantity
            }))
        };
        
        const response = await fetch('/api/orders', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(orderData)
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Ошибка при оформлении заказа: ${errorText}`);
        }
        
        const result = await response.json();
        
        // Очищаем корзину
        state.cart = [];
        localStorage.setItem('cart', JSON.stringify(state.cart));
        updateCart();
        
        // Обновляем список продуктов (чтобы обновить доступное количество)
        fetchProducts();
        
        // Обновляем список заказов
        fetchOrders();
        
        // Показываем сообщение об успешном оформлении заказа
        alert(`Заказ успешно оформлен! Номер заказа: ${result.order_id}`);
        
    } catch (error) {
        console.error('Ошибка при оформлении заказа:', error);
        alert(error.message);
    }
}

// Получение списка заказов
async function fetchOrders() {
    const userId = DOM.userIdFilter.value.trim();
    if (!userId) return;
    
    try {
        const url = new URL('/api/orders', window.location.origin);
        url.searchParams.append('user_id', userId);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Ошибка при получении заказов: ${response.statusText}`);
        }
        
        const orders = await response.json();
        state.orders = orders;
        
        renderOrders();
    } catch (error) {
        console.error('Ошибка при загрузке заказов:', error);
        DOM.ordersList.innerHTML = '<div class="main__orders-error">Ошибка при загрузке заказов</div>';
    }
}

// Отображение списка заказов
function renderOrders() {
    DOM.ordersList.innerHTML = '';
    
    if (state.orders.length === 0) {
        DOM.ordersList.innerHTML = '<div class="main__orders-empty">У вас пока нет заказов</div>';
        return;
    }
    
    state.orders.forEach(order => {
        const orderElement = document.createElement('div');
        orderElement.className = 'main__order-item';
        
        const createdDate = new Date(order.created_at).toLocaleString();
        
        orderElement.innerHTML = `
            <div class="main__order-item-header">
                <div class="main__order-item-date">${createdDate}</div>
                <div class="main__order-item-id">ID: ${order.id}</div>
                <div class="main__order-item-status ${order.status}">${getStatusText(order.status)}</div>
            </div>
            ${order.items ? renderOrderItems(order.items) : '<div class="main__order-item-products-empty">Детали заказа недоступны</div>'}
            <div class="main__order-item-total">
                <span>Итого:</span>
                <span>${order.total_price.toFixed(2)} ₸</span>
            </div>
            ${renderOrderActions(order)}
        `;
        
        DOM.ordersList.appendChild(orderElement);
        
        // Добавляем обработчики для кнопок
        const actionButtons = orderElement.querySelectorAll('.main__order-item-button');
        actionButtons.forEach(button => {
            button.addEventListener('click', () => updateOrderStatus(order.id, button.dataset.status));
        });
    });
}

// Отображение элементов заказа
function renderOrderItems(items) {
    if (!items || items.length === 0) return '';
    
    let html = '<div class="main__order-item-products">';
    
    items.forEach(item => {
        const productName = item.product ? item.product.Name : 'Товар недоступен';
        html += `
            <div class="main__order-item-product">
                <div class="main__order-item-product-name">${productName}</div>
                <div class="main__order-item-product-quantity">${item.quantity} шт.</div>
                <div class="main__order-item-product-price">${item.price.toFixed(2)} ₸</div>
            </div>
        `;
    });
    
    html += '</div>';
    return html;
}

// Отображение кнопок действий для заказа
function renderOrderActions(order) {
    if (order.status === 'completed' || order.status === 'cancelled') {
        return '';
    }
    
    return `
        <div class="main__order-item-actions">
            ${order.status === 'pending' ? `
                <button class="main__order-item-button complete" data-status="completed">Завершить</button>
                <button class="main__order-item-button cancel" data-status="cancelled">Отменить</button>
            ` : ''}
        </div>
    `;
}

// Обновление статуса заказа
async function updateOrderStatus(orderId, status) {
    try {
        const response = await fetch(`/api/orders/${orderId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ status })
        });
        
        if (!response.ok) {
            throw new Error(`Ошибка при обновлении статуса заказа: ${response.statusText}`);
        }
        
        // Обновляем список заказов
        fetchOrders();
        
        // Если заказ был отменен, обновляем список продуктов (чтобы обновить доступное количество)
        if (status === 'cancelled') {
            fetchProducts();
        }
    } catch (error) {
        console.error('Ошибка при обновлении статуса заказа:', error);
        alert(error.message);
    }
}

// Получение текстового представления статуса заказа
function getStatusText(status) {
    switch (status) {
        case 'pending':
            return 'В обработке';
        case 'completed':
            return 'Выполнен';
        case 'cancelled':
            return 'Отменен';
        default:
            return status;
    }
}