window.bus = new Vue();

Vue.component('group-list', {
	template: '#group-list-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			groups: [],
		}
	},
	created: function() {
		this.fetchGroups();
	},
	mounted: function() {
		window.bus.$on('group-created', this.appendGroup);
	},
	methods: {
		appendGroup: function(group) {
			this.groups.push(group);
		},
		fetchGroups: function() {
			axios.get("http://127.0.0.1:5100/api/groups")
			.then(function(resp) {
				if (resp.data.success) {
					this.groups = resp.data.data;
				}
			}.bind(this));
		},
		createNewGroup: function() {
			window.bus.$emit('group-create-modal');
		},
		editGroup: function(group) {
			window.bus.$emit('group-edit-modal', group);
		}
	}
});

Vue.component('group-create-modal', {
	template: '#group-create-modal-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			isActive: false,
			group: {},
		}
	},
	mounted: function() {
		bus.$on('group-create-modal', this.show);
		bus.$on('group-create-modal-close', this.close);
		document.addEventListener("keydown", this.esc);
	},
	beforeDestroy: function() {
		bus.$off('group-create-modal', this.show);
		bus.$off('group-create-modal-close', this.close);
		document.removeEventListener("keydown", this.esc);
	},
	methods: {
		show: function() {
			this.isActive = true;
		},
		close: function() {
			this.isActive = false;
		},
		esc: function(e) {
			if (this.isActive && e.keyCode === 27) {
				this.close();
			}
		},
		send: function(e) {
			e.preventDefault();
			axios.post("http://127.0.0.1:5100/api/groups/create",
				JSON.stringify(this.group),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('group-created', resp.data.data);
					this.close();
				}
			}.bind(this));
			this.reset();
		},
		reset: function() {
			this.group = {};
		}
	}
});

Vue.component('group-edit-modal', {
	template: '#group-edit-modal-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			isActive: false,
			group: {},
		}
	},
	mounted: function() {
		bus.$on('group-edit-modal', this.show);
		bus.$on('group-edit-modal-close', this.close);
		document.addEventListener("keydown", this.esc);
	},
	beforeDestroy: function() {
		bus.$off('group-edit-modal', this.show);
		bus.$off('group-edit-modal-close', this.close);
		document.removeEventListener("keydown", this.esc);
	},
	methods: {
		show: function(group) {
			this.isActive = true;
			this.group = group;
		},
		close: function() {
			this.isActive = false;
		},
		esc: function(e) {
			if (this.isActive && e.keyCode === 27) {
				this.close();
			}
		},
		send: function(e) {
			e.preventDefault();
			axios.post("http://127.0.0.1:5100/api/groups/" + this.group.id + "/edit",
				JSON.stringify(this.group),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('group-updated', resp.data.data);
					this.close();
				}
			}.bind(this));
			this.reset();
		},
		reset: function() {
			this.group = {};
		}
	}
});

Vue.component('task-list', {
	template: '#task-list-template',
	delimiters: ['${', '}'],
	props: ['group'],
	data: function() {
		return {
			tasks: [],
		}
	},
	created: function() {
		this.fetchTasks(this.group.id);
	},
	mounted: function() {
		window.bus.$on('task-created', this.appendTask);
	},
	methods: {
		fetchTasks:function(groupid) {
			axios.get("http://127.0.0.1:5100/api/tasks?groupid=" + groupid)
			.then(function(resp) {
				if (resp.data.success) {
					this.tasks = resp.data.data;
				}
			}.bind(this));
		},
		appendTask: function(task) {
			if (task.groupid === this.group.id) {
				this.tasks.push(task);
			}
		},
		createNewTask: function() {
			window.bus.$emit('task-create-modal', this.group);
		},
		editTask: function(task) {
			window.bus.$emit('task-edit-modal', task);
		},
		toggleTask: function(task) {
			task.completed = !task.completed;
			if (task.completed) {
				window.bus.$emit('task-completed', task);
			}
			axios.post("http://127.0.0.1:5100/api/tasks/" + task.id + "/edit",
				JSON.stringify(task),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('task-updated', resp.data.data);
				}
			}.bind(this));
		}
	}
});

Vue.component('task-create-modal', {
	template: '#task-create-modal-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			isActive: false,
			task: {},
			group: {},
		}
	},
	mounted: function() {
		bus.$on('task-create-modal', this.show);
		bus.$on('task-create-modal-close', this.close);
		document.addEventListener("keydown", this.esc);
	},
	beforeDestroy: function() {
		bus.$off('task-create-modal', this.show);
		bus.$off('task-create-modal-close', this.close);
		document.removeEventListener("keydown", this.esc);
	},
	methods: {
		show: function(group) {
			this.isActive = true;
			this.group = group;
		},
		close: function() {
			this.isActive = false;
		},
		esc: function(e) {
			if (this.isActive && e.keyCode === 27) {
				this.close();
			}
		},
		send: function(e) {
			e.preventDefault();
			axios.post("http://127.0.0.1:5100/api/tasks/create",
				JSON.stringify({
					groupid: this.group.id,
					label: this.task.label,
					completed: this.task.completed
				}),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('task-created', resp.data.data);
					this.close();
				}
			}.bind(this));
			this.reset();
		},
		reset: function() {
			this.task = {};
		}
	}
});

Vue.component('task-edit-modal', {
	template: '#task-edit-modal-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			isActive: false,
			task: {},
		}
	},
	mounted: function() {
		bus.$on('task-edit-modal', this.show);
		bus.$on('task-edit-modal-close', this.close);
		document.addEventListener("keydown", this.esc);
	},
	beforeDestroy: function() {
		bus.$off('task-edit-modal', this.show);
		bus.$off('task-edit-modal-close', this.close);
		document.removeEventListener("keydown", this.esc);
	},
	methods: {
		show: function(task) {
			this.isActive = true;
			this.task = task;
		},
		close: function() {
			this.isActive = false;
		},
		esc: function(e) {
			if (this.isActive && e.keyCode === 27) {
				this.close();
			}
		},
		send: function(e) {
			e.preventDefault();
			axios.post("http://127.0.0.1:5100/api/tasks/" + this.task.id + "/edit",
				JSON.stringify(this.task),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('task-updated', resp.data.data);
					this.close();
				}
			}.bind(this));
			this.reset();
		},
		reset: function() {
			this.task = {};
		}
	}
});

Vue.component('login-modal', {
	template: '#login-modal-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			isActive: false,
			creds : {}
		}
	},
	mounted: function() {
		bus.$on('login-modal', this.show);
		bus.$on('login-modal-close', this.close);
		document.addEventListener("keydown", this.esc);
	},
	beforeDestroy: function() {
		bus.$off('login-modal', this.show);
		bus.$off('login-modal-close', this.close);
		document.removeEventListener("keydown", this.esc);
	},
	methods: {
		show: function() {
			this.isActive = true;
		},
		close: function() {
			this.isActive = false;
		},
		esc: function(e) {
			if (this.isActive && e.keyCode === 27) {
				this.close();
			}
		},
		send: function(e) {
			e.preventDefault();
			axios.post("http://127.0.0.1:5100/api/auth/login",
				JSON.stringify(this.creds),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('logged-in', resp.data.data);
					this.close();
				}
			}.bind(this));
			this.reset();
		},
		reset: function() {
			this.creds = {};
		}
	}
});

// pages
var homepage = {
	template: '#homepage-template'
};

var tasks = {
	template: '#taskspage-template'
};

var groups = {
	template: '#groupspage-template'
}

// vue vm
var vm = new Vue({
	el: '#app',
	delimiters: ['${', '}'],
	data: {
		isNavOpen: false,
		currentView: 'homepage',
	},
	components: {
		homepage: homepage,
		tasks: tasks,
		groups: groups
	},
	methods: {
		toggleNav: function() {
			this.isNavOpen = !this.isNavOpen;
		},
		load: function(page) {
			this.currentView = page;
		},
		loginModal: function() {
			window.bus.$emit('login-modal')
		}
	}
});