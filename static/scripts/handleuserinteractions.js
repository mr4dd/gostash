const addButtonPressed = document.querySelector(".new");
const InfoContainer = document.querySelector(".infoContainer");
const form = document.querySelector(".formData");
const entries = document.querySelector(".entries");
let isVisible = false;

document.body.addEventListener('keydown', function(e) {
  if (e.key == "Escape") {
  	isVisible = true;
  	ContainerDisplayManager();
  }
});

function ContainerDisplayManager() {
	if (isVisible){
		InfoContainer.style.display = "none";
		isVisible = false;
	}else{
		InfoContainer.style.display = "flex";
		isVisible = true;
	}	
};

function newEntry(){
	form.children[4].value = 0;
	form.children[0].value = parseInt(entries.children[entries.children.length - 1].getAttribute("data-id")) + 1 || 0;
	ContainerDisplayManager();
}

const editButtonPressed = (element) => {
	const form = document.querySelector(".formData");
	form.children[4].value = 1; //for endpoint selection
	const id = form.children[0]
	const firstParent = element.parentNode.parentNode;
	const quantityValue = firstParent.querySelector('.count').innerText;
	const quantity = form.children[1]
	const nameValue = firstParent.querySelector('.titleh2').innerText;
	const name = form.children[2]
	const categoryValue = firstParent.querySelector('.tag').innerText;
	const category = form.children[3]
	const descriptionValue = firstParent.querySelector('.description').innerText;
	const description = form.children[5]
	
	id.value = firstParent.getAttribute("data-id");
	quantity.value = quantityValue;
	name.value = nameValue;
	category.value = categoryValue;
	description.value = descriptionValue;
	
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
  		if (value == 0) {
		  	endpoint = '/additem';
		  } else {
		  	endpoint = '/edititem';
		  }
  	}else if (key == "id" || key == "quantity") {
  		jsonObject[key] = parseInt(value);	
  	}
  	else {
    	jsonObject[key] = value;
    }
  });

  fetch(endpoint, {
  	method: 'POST',
  	body: JSON.stringify(jsonObject)
  })
  clearForm()
  ContainerDisplayManager();
  fetchDataAndRender();
});

function remove(element) {
		const confirmation = confirm("are you sure you want to delete this entry?");
		if (confirmation == false) {
			return;
		} else {
			const parent = element.parentNode.parentNode;
			const jsonObject = {}
			jsonObject["id"] = parseInt(parent.getAttribute("data-id"));
			console.log(jsonObject["id"]);
			jsonObject["quantity"] = 0;
			jsonObject["name"] = "null";
			jsonObject["category"] = "null";
			fetch("/deleteitem", {
				method: 'POST',
				body: JSON.stringify(jsonObject)
			});

			parent.remove()
	 }
} 