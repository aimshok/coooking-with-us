let recipes = [];
let totalCount = 0; // Total number of Pokémon
let currentPage = 1;
let recipesPerPage = 6;
let sortAttribute = "name";
let sortDirection = "asc";
let levelFilter = "All";

async function fetchRecipes() {
    try {
        const response = await fetch(`/recipes?page=${currentPage}&perPage=${recipesPerPage}&sortBy=${sortAttribute}&sortDirection=${sortDirection}&levelFilter=${levelFilter}&minPrice=${minPrice}&maxPrice=${maxPrice}`);

        if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

        const data = await response.json();
        console.log("Received data:", data);

        if (!data || !Array.isArray(data.Recipes)) {
            throw new Error("Invalid response format: Missing 'Recipes' array");
        }

        recipes = data.Recipes;
        totalCount = data.totalCount;
        displayRecipes();
        displayPagination();
    } catch (error) {
        console.error("Failed to fetch recipes:", error);
    }
}


function showNumrecipes(num) {
    recipesPerPage = num;
    currentPage = 1;
    fetchRecipes();
}


function filterByLevel(level) {
    levelFilter = level;
    currentPage = 1;
    fetchRecipes();
}

function sortrecipes(attribute, direction) {
    sortAttribute = attribute;
    sortDirection = direction;
    currentPage = 1;
    fetchRecipes();
}

const cart = [];

function toggleCart() {
    const cartModal = document.getElementById("cart");
    cartModal.style.display = cartModal.style.display === "block" ? "none" : "block";
    displayCartItems();
}

function closeInfo() {
    document.getElementById("recipeInfo").style.display = "none";
}

function viewMore(recipe) {
    const recipeInfo = document.getElementById("recipeDetails");
    recipeInfo.innerHTML = `
            <h3>${recipe.name}</h3>
            <img src="${recipe.path}" alt="${recipe.name}" style="max-width: 140px;">
            <p><strong>Level:</strong> ${recipe.level}</p>
            <p><strong>Calories:</strong>${recipe.calories}</p>
            <p>${recipe.desc}</p>
            <p><strong>Price:</strong> $${recipe.price}</p>
            <button onclick="addToCart(${JSON.stringify(recipe).replace(/"/g, '&quot;')})">Buy</button>
        `;
    document.getElementById("recipeInfo").style.display = "block";
}

function addToCart(recipe) {
    cart.push(recipe);
    alert(`${recipe.name} has been added to the cart!`);
    closeInfo()
}

function displayCartItems() {
    const cartItemsContainer = document.getElementById("cartItems");
    cartItemsContainer.innerHTML = "";

    if (cart.length === 0) {
        cartItemsContainer.innerHTML = "<p>Your cart is empty.</p>";
        return;
    }

    // Total price calculation
    let totalPrice = 0;
    cart.forEach((item, index) => {
        const cartItem = document.createElement("div");
        cartItem.innerHTML = `
                <h4>${item.name}</h4>
                <p>Level: ${item.level}</p>
                <p>Price: $${item.price}</p>
                <button onclick="removeFromCart(${index})">Remove</button>
            `;
        cartItemsContainer.appendChild(cartItem);

        totalPrice += item.price; // Add price of each item to totalPrice
    });

    // Display the total price
    const totalPriceContainer = document.createElement("div");
    totalPriceContainer.innerHTML = `<h3>Total Price: $${totalPrice.toFixed(2)}</h3>`;
    cartItemsContainer.appendChild(totalPriceContainer);
}

function removeFromCart(index) {
    cart.splice(index, 1);
    displayCartItems();
}

function checkout() {
    if (cart.length === 0) {
        alert("Your cart is empty!");
        return;
    }

    alert("Thank you for your purchase!");
    cart.length = 0;
    toggleCart();
}

function displayRecipes() {
    console.log("Recipes:", recipes);  // Log the recipes array to see its value
    // Check if recipes is an array and if it's not empty
    if (!Array.isArray(recipes) || recipes.length === 0) {
        console.log("No recipes to display");
        return; // No recipes to display, exit the function
    }

    const container = document.querySelector(".recipe-container");
    container.innerHTML = "";

    recipes.forEach(recipe => {
        const card = document.createElement("div");
        console.log("Image Path: ", recipe.imagePath);
        card.className = "recipe-card";
        card.innerHTML = `
            <h3>${recipe.name}</h3>
            <img src="${recipe.path}" alt="${recipe.name}">
            <p><strong>Level:</strong> ${recipe.level}</p>
            <p><strong>Calories:</strong>${recipe.calories}</p>
            <p><strong>Price:</strong> $${recipe.price}</p>
            <button onclick="viewMore(${JSON.stringify(recipe).replace(/"/g, '&quot;')})">View More</button>
        `;
        container.appendChild(card);
    });
}
let minPrice = 0;
let maxPrice = 1000;

function filterByPrice() {
    minPrice = document.getElementById("minPrice").value;
    maxPrice = document.getElementById("maxPrice").value;

    // Обновляем отображаемые значения цен
    document.getElementById("minPriceValue").innerText = minPrice;
    document.getElementById("maxPriceValue").innerText = maxPrice;
    const minPriceSlider = document.getElementById("minPrice");
    minPriceSlider.max = maxPrice;

    // Устанавливаем страницу на первую, чтобы начать с первой страницы результатов
    currentPage = 1;
    fetchRecipes();
}
document.getElementById("maxPrice").addEventListener("input", function() {
    const maxValue = this.value;
    document.getElementById("minPrice").max = maxValue; // Устанавливаем max для minPrice равным maxPrice
});

function displayPagination() {
    const paginationContainer = document.querySelector(".pagination");
    paginationContainer.innerHTML = "";

    const totalPages = Math.ceil(totalCount / recipesPerPage);

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement("button");
        button.innerText = i;
        button.onclick = () => {
            currentPage = i;
            fetchRecipes();
        };
        paginationContainer.appendChild(button);
    }
}

function populateFilter() {
    const filterContainer = document.querySelector("#recipeLevelFilter");
    const level = ["All", "beginner", "amateur","master"];
    level.forEach(level => {
        const option = document.createElement("option");
        option.value = level;
        option.innerText = level;
        filterContainer.appendChild(option);
    });
}
function numPerPage() {
    const numContainer = document.querySelector("#recipeNum");
    const nums = [6, 12, 18, 24];
    nums.forEach(num => {
        const option = document.createElement("option");
        option.value = num.toString();
        option.innerText = num.toString();
        numContainer.appendChild(option);
    });
}
window.onload = () => {
    populateFilter();
    numPerPage();
    fetchRecipes();
};