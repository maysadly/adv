document
  .getElementById("addProduct-form")
  .addEventListener("submit", async function (event) {
    event.preventDefault();

    const productName = document.querySelector(
      'input[name="productName"]'
    ).value;
    const productPrice = parseFloat(
      document.querySelector('input[name="productPrice"]').value
    );
    const productStock = parseInt(
      document.querySelector('input[name="productStock"]').value
    );

    try {
      const response = await fetch("/api/products", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          Name: productName,
          Price: productPrice,
          Stock: productStock,
        }),
      });

      if (response.ok) {
        document.getElementById("addProduct-form").reset();
        fetchProducts();
      } else {
        const error = await response.text();
        alert(`Failed to add product: ${error}`);
      }
    } catch (error) {
      alert(`Failed to add product: ${error}`);
    }
  });

async function fetchProducts() {
  try {
    const response = await fetch("/api/products");
    if (!response.ok) {
      throw new Error("Failed to fetch products");
    }
    const products = await response.json();

    const productsList = document.getElementById("products-list");
    productsList.innerHTML = "";

    products.forEach((product) => {
      const productItem = document.createElement("div");
      productItem.className = "main__products-item";
      productItem.innerHTML = `
        <div class="main__products-item-wrap">
          <div class="main__products-item-name">${product.Name}</div>
          <button class="main__products-item-edit" onclick="editProduct('${product.ID}', '${product.Name}', ${
        product.Price
      }, ${product.Stock})">✏️</button>

      
          <button class="main__products-item-delete" onclick="deleteProduct('${
            product.ID
          }')">❌</button>
        </div>
        <div class="main__products-item-price">${product.Price.toFixed(
          2
        )}₸</div>
        <div class="main__products-item-stock">Stock: ${product.Stock}</div>
      `;
      productsList.appendChild(productItem);
    });
  } catch (error) {
    console.error("Error fetching products:", error);
  }
}

async function deleteProduct(productId) {
  try {
    const response = await fetch(`/api/products/${productId}`, {
      method: "DELETE",
    });
    if (response.ok) {
      fetchProducts();
    } else {
      console.error("Error deleting product:", response.statusText);
    }
  } catch (error) {
    console.error("Error deleting product:", error);
  }
}

function editProduct(id, name, price, stock) {
  const productItem = document
    .querySelector(`button[onclick*="editProduct('${id}'"]`)
    .closest(".main__products-item");

  productItem.innerHTML = `
      <form class="main__products-item-update" onsubmit="submitUpdate(event, '${id}')">
        <input class="main__products-item-update-input" type="text" name="Name" value="${name}" required />
        <input class="main__products-item-update-input" type="number" name="Price" step="0.01" value="${price}" required />
        <input class="main__products-item-update-input" type="number" name="Stock" value="${stock}" required />

        <div class="main__products-item-update-buttons"> 
            <button class="main__products-item-update-buttons-save" type="submit">Save</button>
            <button class="main__products-item-update-buttons-cancel" type="button" onclick="fetchProducts()">Cancel</button>
        </div>
      </form>
    `;
}

async function submitUpdate(event, id) {
  event.preventDefault();

  const form = event.target;
  const updatedProduct = {
    Name: form.Name.value,
    Price: parseFloat(form.Price.value),
    Stock: parseInt(form.Stock.value),
  };

  try {
    const response = await fetch(`/api/products/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(updatedProduct),
    });

    if (response.ok) {
      fetchProducts();
    } else {
      const error = await response.text();
      alert(`Failed to update product: ${error}`);
    }
  } catch (error) {
    alert(`Failed to update product: ${error}`);
  }
}

async function deleteProduct(productId) {
  try {
    const response = await fetch(`/api/products/${productId}`, {
      method: "DELETE",
    });
    if (response.ok) {
      fetchProducts();
    } else {
      console.error("Error deleting product:", response.statusText);
    }
  } catch (error) {
    console.error("Error deleting product:", error);
  }
}

window.onload = fetchProducts;
