document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById('addRecipeForm');
    form.reset();

    const addRecipeButton = document.getElementById("addRecipeButton");
    addRecipeButton.addEventListener("click", addRecipe);
});


async function addRecipe(event) {
    event.preventDefault();  // Prevent the form from submitting
    const form = document.getElementById('addRecipeForm');

    const name = document.getElementById("RecipeName").value;
    const level = document.getElementById("RecipeLevel").value;
    const calories = document.getElementById("RecipeCalorie").value;
    const description = document.getElementById("RecipeDescription").value;
    const imagePath = document.getElementById("RecipeImage").files[0];
    const price = document.getElementById("RecipePrice").value;

    if (!name.trim() || !description.trim() || !imagePath) {
        alert("All fields are required.");
        return;
    }

    const formData = new FormData();
    formData.append("name", name);
    formData.append("level", level);
    formData.append("calories", calories);
    formData.append("price", price);
    formData.append("description", description);
    formData.append("image", imagePath);  // Change this to 'image' to match server side

    const response = await fetch(`/addRecipe`, {
        method: "POST",
        body: formData,  // Don't set the Content-Type manually
    });
    form.reset();
    const result = await response.json();
    if (response.ok) {
        alert("Recipe added successfully!");
    } else {
        alert("Failed to add Recipe: " + result.message);
    }
}

    async function grantAdmin() {
        const email = document.getElementById("userEmail").value;

        if (!email.trim()) {
            alert("Email is required.");
            return;
        }

        const response = await fetch(`/grantAdmin`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({email})
        });

        const result = await response.json();
        if (response.ok) {
            alert("User granted admin rights!");
        } else {
            alert("Failed to grant admin rights: " + result.message);
        }
    }

    async function deleteRecipe() {
        const RecipeId = document.getElementById("RecipeId").value;

        if (!RecipeId.trim()) {
            alert("Recipe ID is required.");
            return;
        }

        const response = await fetch(`/deleteRecipe`, {
            method: "DELETE",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({id: RecipeId})
        });

        const result = await response.json();
        if (response.ok) {
            alert("Recipe deleted successfully!");
        } else {
            alert("Failed to delete Recipe: " + result.message);
        }
    }
