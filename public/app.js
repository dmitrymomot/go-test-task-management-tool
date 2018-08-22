axios.defaults.headers.common['Content-Type'] = 'application/json';
axios.defaults.headers.common['Accept'] = 'application/json';

var app = new Vue({
  el: '#app',

  created() {
    this.reload()
  },

  data() {
    return {
      message: null,
      allTasks: null,
      newTasks: null,
      completedTasks: null,
      createForm: {
        error: null,
        title: null,
        description: null,
      },
      editForm: {
        error: null,
        title: null,
        description: null,
        done: false,
        status: null,
      },
      editingTask: {
        id: null,
        title: null,
        description: null,
        status: null,
        created_at: null,
        completed_at: null,
      },
    }
  },

	methods: {
		reload() {
			this.getAll()
			this.getNew()
			this.getCompleted()
		},
		getAll() {
      axios.get('/api/tasks')
        .then(response => {
          console.log(response.data)
          this.allTasks = response.data.data;
        });
		},
		getNew() {
      axios.get('/api/tasks/new')
        .then(response => {
          console.log(response.data)
          this.newTasks = response.data.data;
        });
		},
		getCompleted() {
      axios.get('/api/tasks/completed')
        .then(response => {
          console.log(response.data)
          this.completedTasks = response.data.data;
        });
		},
		createTask() {
      axios.post('/api/tasks', this.createForm)
        .then(response => {
          this.message = 'New task has been created';
          this.$nextTick(function () {
            this.createForm.error = null
            this.createForm.title = null
            this.createForm.description = null
            this.reload();
          })
        })
        .catch(error => {
          this.createForm.error = error.response.data.error
        });
		},
		editTask(item) {
      this.editingTask = item
      this.editForm.title = this.editingTask.title
      this.editForm.description = this.editingTask.description
      this.editForm.done = this.editingTask.status == "done" ? true : false;
      $('#editModal').modal('show')
		},
		updateTask() {
      console.log(this.editForm)
      this.editForm.status = this.editForm.done ? "done" : "new";
      axios.put('/api/tasks/'+this.editingTask.id, this.editForm)
        .then(response => {
          this.$nextTick(function () {
            this.reload();
            $('#editModal').modal('hide')
          })
        })
        .catch(error => {
          this.editForm.error = error.response.data.error
        });
		},
		doneTask(item) {
      axios.patch('/api/tasks/'+item.id+'/done', {data:{}})
        .then(response => {
          this.tasks = response.data.data;
          this.$nextTick(function () {
              this.reload();
          })
        });
		},
		deleteTask(item) {
      axios.delete('/api/tasks/'+item.id, {data:{}})
        .then(response => {
          this.tasks = response.data.data;
          this.$nextTick(function () {
              this.reload();
          })
        });
		},
    hasError(error, key) {
      return (error != null && error[key])
    },
    getError(error, key) {
      if (this.hasError(error, key)) {
        return error[key]
      }
      return null
    }
	}
})

$('.popover-dismiss').popover({
  trigger: 'focus'
})

$('.collapse').collapse()
