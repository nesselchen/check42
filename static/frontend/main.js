function query(selector) {
    return document.querySelector(selector)
}

function elem(name) {
    return (params) => {
        const el = document.createElement(name)
        if (params.id) {
            el.id = params.id 
        }
        if (params.class) {
            el.className = params.class
        }
        if (params.text) {
            text = document.createTextNode(params.text)
            el.appendChild(text)
        }
        return el
    }
}

const li = elem("li")
const div = elem("div")
const btn = elem("button")
const span = elem("span")

const DeleteButton = (id, parent) => {
    const button = btn({
        text: "Delete",
        class: "delete-btn"
    })
    button.addEventListener("click", async (e) => {
        const deleted = await deleteTodo(id)
        if (deleted) {
            parent.remove()
        }
    })
    return button
}

const TodoComponent = (id, text, done) => { 
    const todo = li({
        class: "todo" + (done ? " done" : ""),
    })
    const check = btn({
        text: done ? "Do it again": "Done",
    })
    check.addEventListener("click", event => {
        toggleTodo(id, !done).then(ok => {
            if (ok) {
                todo.replaceWith(TodoComponent(id, text, !done))
            }
        })
    })
    controls = div({
        class: "controls",
    })
    todo.appendChild(span({text, class: "todo-text"}))
    controls.appendChild(check)
    controls.appendChild(DeleteButton(id, todo))
    todo.appendChild(controls)
    return todo
}

async function initializePage() {
    let res = await fetch("/api/todo")
    if (res.status == 401) {
        const username = prompt("You're not logged in. What's your username?")
        const password = prompt("And now your password?")
        await loginUser(username, password)

        res = await fetch("/api/todo")
        if (res.status == 401) {
            alert("Something's not right. Try refreshing the page.")
            return
        }
    }
    const todos = await res.json()
    registerTodos(todos)
}

function registerTodos(todos) {
    todos.forEach(t => {
        const list = query("#todos")
        const item = TodoComponent(t.id, t.text, t.done)
        list.appendChild(item)
    })
}

async function toggleTodo(id, newVal) {
    const res = await fetch(`/api/todo/${id}?done=${newVal}`, {
        method: "PATCH"
    })
    return res.ok
}

async function deleteTodo(id) {
    const res = await fetch(`/api/todo/${id}`, {
        method: "DELETE"
    })
    return res.ok
}

// TODO: how does a user login?
async function loginUser(username, password) {
    const encoded = btoa(username + ":" + password) 
    const res = await fetch("/auth/login", {
        method: "POST",
        headers: {
            Authorization: `Basic ${encoded}`,
        },
    })
    return res.status < 400
}

query(".todo-form").addEventListener("submit", e => {
    e.preventDefault()
    const data = new FormData(e.target)
    const todo = Object.fromEntries(data)
    fetch("/api/todo", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(todo),
    })
    .then(res => {
        if (!res.statusCode == 201) {
            throw("Error creating todo")    
        }
        return res.json()
    })
    .then(id => {
        query("#todos").appendChild(
            TodoComponent(id, todo.text, false)
        )
        query(".todo-form").reset()
    })
    .catch(e => {
        console.log(e)
    })

}, false)

initializePage()