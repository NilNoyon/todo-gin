async function fetchTodos() {
     const res = await fetch("/api/todos");
     const todos = await res.json();
     const list = document.getElementById("todoList");
     list.innerHTML = "";
     todos.forEach(todo => {
       const li = document.createElement("li");
       li.textContent = todo.title;
       if (todo.done) li.classList.add("done");
   
       const doneBtn = document.createElement("button");
       doneBtn.textContent = "✓";
       doneBtn.onclick = () => markDone(todo.id);
   
       const delBtn = document.createElement("button");
       delBtn.textContent = "✕";
       delBtn.onclick = () => deleteTodo(todo.id);
   
       li.appendChild(doneBtn);
       li.appendChild(delBtn);
       list.appendChild(li);
     });
   }
   
   async function addTodo() {
     const title = document.getElementById("newTodo").value;
     if (!title) return;
     await fetch("/api/todos", {
       method: "POST",
       headers: { "Content-Type": "application/json" },
       body: JSON.stringify({ title })
     });
     document.getElementById("newTodo").value = "";
     fetchTodos();
   }
   
   async function markDone(id) {
     await fetch(`/api/todos/${id}`, { method: "PUT" });
     fetchTodos();
   }
   
   async function deleteTodo(id) {
     await fetch(`/api/todos/${id}`, { method: "DELETE" });
     fetchTodos();
   }
   
   fetchTodos();
   