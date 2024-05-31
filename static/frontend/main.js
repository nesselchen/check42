function query(selector) {
    return document.querySelector(selector)
}

function elem(name) {
    return (params) => {
        const el = document.createElement(name)
        if (params.id != undefined) {
            el.id = params.id
        }
        if (params.class != undefined) {
            el.className = params.class
        }
        if (params.text != undefined) {
            text = document.createTextNode(params.text)
            el.appendChild(text)
        }
        if (params.value) {
            el.value = params.value
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
        text: done ? "Do it again" : "Done",
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
    todo.appendChild(span({ text, class: "todo-text" }))
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
    loadCategories()
    const todos = await res.json()
    registerTodos(todos)
}

async function loadCategories() {
    const res = await fetch("/api/todo/category")
    const parsed = await res.json()
    categories.clear()
    categories.set("My todos", {})
    const dd = query("#dropdown-category")
    parsed.forEach(cat => {
        console.log(cat)
        categories.set(cat.name, cat)
        dd.appendChild(elem("option")({
            value: cat.name,
            text: cat.name,
        }))
    })
}

function registerTodos(todos) {
    const categories = Object.groupBy(todos, t => t.category.name || "My todos")
    for (const cat in categories) {
        const list = query("#categories")
        const catList = div({ text: cat, class: "category", id: cat})
        list.appendChild(catList)
        const todoList = document.createElement("ul")
        catList.appendChild(todoList)
        categories[cat].forEach(t => {
            const item = TodoComponent(t.id, t.text, t.done)
            todoList.appendChild(item)
        })
    }
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

query(".todo-form")
    .addEventListener("submit", e => {
        e.preventDefault()
        const data = new FormData(e.target)
        let todo = Object.fromEntries(data)
        const cat = query("#dropdown-category").value
        todo.category = categories.get(cat)
        
        fetch("/api/todo", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(todo),
        })
            .then(res => {
                if (!res.statusCode == 201) {
                    throw ("Error creating todo")
                }
                return res.json()
            })
            .then(id => {
                console.log(cat)
                const catDiv = document.getElementById(cat)
                catDiv.appendChild(
                    TodoComponent(id, todo.text, false)
                )
                query(".todo-form").reset()
            })
            .catch(e => {
                console.log(e)
            })

    }, false)

    
const categories = new Map()

initializePage()
console.log("Initialized")