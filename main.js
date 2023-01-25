class Item {
  constructor(id, age, name) {
    this.userId = id;
    this.userAge = age;
    this.userName = name;
    this.itemElement = document.createElement("div");
    this.itemElement.innerText = `Id: ${this.userId} | Age: ${this.userAge} | Name: ${this.userName}`;
  }

  getItemElement() {
    return this.itemElement;
  }

  removeItemElement() {
    this.itemElement.remove();
  }
}

class DeleteButton {
  constructor(linkedItem) {
    this.buttonElement = document.createElement("button");
    this.buttonElement.innerText = "X";
    this.buttonElement.addEventListener(
      "click",
      this.onClickHandler.bind(this)
    );
    this.linkedItem = linkedItem;
  }

  onClickHandler(event) {
    this.remove();
  }

  getButtonElement() {
    return this.buttonElement;
  }

  remove() {
    this.linkedItem.removeItemElement();
    this.buttonElement.remove();
  }
}

async function load(fetchUrl) {
  const obj = await fetchItems(fetchUrl);
  displayItems(obj);
}

async function fetchItems(fetchUrl) {
  const response = await fetch(fetchUrl);
  if (response.status !== 200) {
    alert("Response error");
  }
  const obj = await response.json();
  return obj;
}

function displayItems(obj) {
  const subContainer = document.createElement("div");
  const item = new Item(obj.UserId, obj.UserAge, obj.UserName);
  const button = new DeleteButton(item);
  subContainer.appendChild(item.getItemElement());
  subContainer.appendChild(button.getButtonElement());
  itemContainer.appendChild(subContainer);
}

const itemContainer = document.getElementById("item_container");
const fetchUrl = "http://localhost:8000/get";
load(fetchUrl);
