const addButtonPressed = document.querySelector(".new");
const InfoContainer = document.querySelector(".infoContainer");
const form = document.querySelector(".formData");
let isVisible = false;

function ContainerDisplayManager() {
	if (isVisible){
		InfoContainer.style.display = "none";
		isVisible = false;
	}else{
		InfoContainer.style.display = "flex";
		isVisible = true;
	}	
};

addButtonPressed.addEventListener("click", () => {
	form.children[4].value = 1;
	ContainerDisplayManager();
});

const editButtonPressed = (element) => {
	const form = document.querySelector(".formData");
	form.children[4].value = 1;
	const id = form.children[0]
	const firstParent = element.parentNode.parentNode;
	const quantityValue = firstParent.querySelector('.count').innerText;
	const quantity = form.children[1]
	const nameValue = firstParent.querySelector('.entrytitle').innerText;
	const name = form.children[2]
	const categoryValue = firstParent.querySelector('.tag').innerText;
	const category = form.children[3]
	
	id.value = firstParent.getAttribute("data-id");
	quantity.value = quantityValue;
	name.value = nameValue;
	category.value = categoryValue;
	
	ContainerDisplayManager();
};

function clearForm() {
	const children = form.children;
	for (let i = 0; i < children.length - 1; i++){
		children[i].value = '';
	}
}

form.addEventListener("submit", function(event) {
  let endpoint = '';
  event.preventDefault(); // Prevent form submission
  // Get form data
  const formData = new FormData(this);
  
  // Convert form data to JSON object
  const jsonObject = {};
  formData.forEach((value, key) => {
  	if (key == "newselector") {
  		key = "newselector";
  	}else if (key == "id" || key == "quantity") {
  		jsonObject[key] = parseInt(value);	
  	}
  	else {
    	jsonObject[key] = value;
    }
  });
  if (formData['newselector'] == 0) {
  	endpoint = '/additem';
  } else {
  	endpoint = '/edititem';
  }

  fetch(endpoint, {
  	method: 'POST',
  	body: JSON.stringify(jsonObject)
  })
  clearForm()
  ContainerDisplayManager();
  fetchDataAndRender();
});

function remove(element) {
		const parent = element.parentNode;
		const jsonObject = {}
		jsonObject["id"] = parseInt(parent.getAttribute("data-id"));
		jsonObject["quantity"] = 0;
		jsonObject["name"] = "null";
		jsonObject["category"] = "null";
		fetch("/deleteitem", {
			method: 'POST',
			body: JSON.stringify(jsonObject)
		});

		parent.remove()
} 