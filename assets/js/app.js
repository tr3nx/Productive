window.bus = new Vue();

/* Group */
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
		window.bus.$on('group-deleted', this.removeGroup);
	},
	methods: {
		fetchGroups: function() {
			axios.get("http://127.0.0.1:5100/api/groups?userid=" + this.$root.user.id)
			.then(function(resp) {
				if (resp.data.success) {
					this.groups = resp.data.data;
				}
			}.bind(this));
		},
		appendGroup: function(group) {
			this.groups.push(group);
		},
		removeTask: function(task) {
			if (this.group.id === task.groupid) {
				this.tasks = this.tasks.filter(function(t) {
					return (t.id !== task.id);
				});
			}
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
				JSON.stringify({
					userid: vm.user.id,
					order: this.group.order,
					label: this.group.label
				}),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('group-created', resp.data.data);
					this.close();
					this.reset();
				}
			}.bind(this));
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
		remove: function() {
			axios.post("http://127.0.0.1:5100/api/groups/" + this.group.id + "/delete",
				JSON.stringify(this.group),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('group-removed', this.group);
					this.close();
					this.reset();
				}
			}.bind(this));
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
					this.reset();
				}
			}.bind(this));
		},
		reset: function() {
			this.group = {};
		}
	}
});

/* Task */
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
		window.bus.$on('task-deleted', this.removeTask);
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
		removeTask: function(task) {
			if (this.group.id === task.groupid) {
				this.tasks = this.tasks.filter(function(t) {
					return (t.id !== task.id);
				});
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
					userid: vm.user.id,
					label: this.task.label,
					completed: this.task.completed
				}),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('task-created', resp.data.data);
					this.close();
					this.reset();
				}
			}.bind(this));
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
		remove: function() {
			axios.post("http://127.0.0.1:5100/api/tasks/" + this.task.id + "/delete",
				JSON.stringify(this.task),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('task-deleted', resp.data.data);
					this.close();
					this.reset();
				}
			}.bind(this));
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
					this.reset();
				}
			}.bind(this));
		},
		reset: function() {
			this.task = {};
		}
	}
});

/* Auth */
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
					this.reset();
				}
			}.bind(this));
		},
		reset: function() {
			this.creds = {};
		},
		loginModal: function() {
			window.bus.$emit('login-modal');
		},
		signupModal: function() {
			window.bus.$emit('signup-modal');
		}
	}
});

Vue.component('signup-modal', {
	template: '#signup-modal-template',
	delimiters: ['${', '}'],
	data: function() {
		return {
			isActive: false,
			creds : {}
		}
	},
	mounted: function() {
		bus.$on('signup-modal', this.show);
		bus.$on('signup-modal-close', this.close);
		document.addEventListener("keydown", this.esc);
	},
	beforeDestroy: function() {
		bus.$off('signup-modal', this.show);
		bus.$off('signup-modal-close', this.close);
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
			axios.post("http://127.0.0.1:5100/api/auth/signup",
				JSON.stringify(this.creds),
				{ headers : { 'Content-Type' : 'application/json' } })
			.then(function(resp) {
				if (resp.data.success) {
					window.bus.$emit('signed-up', resp.data.data);
					this.close();
					this.reset();
				}
			}.bind(this));
		},
		reset: function() {
			this.creds = {};
		},
		loginModal: function() {
			window.bus.$emit('login-modal');
		},
		signupModal: function() {
			window.bus.$emit('signup-modal');
		}
	}
});

/* Pages */
var homepage = {
	template: '#homepage-template',
	methods: {
		loginModal: function() {
			window.bus.$emit('login-modal');
		},
		signupModal: function() {
			window.bus.$emit('signup-modal');
		}
	}
};

var tasks = {
	template: '#taskspage-template',
	methods: {
		loginModal: function() {
			window.bus.$emit('login-modal');
		},
		signupModal: function() {
			window.bus.$emit('signup-modal');
		}
	}
};

/* Vue vm */
var vm = new Vue({
	el: '#app',
	delimiters: ['${', '}'],
	data: {
		isNavOpen: false,
		currentView: 'homepage',
		isLoggedIn: false,
		user: {},
	},
	components: {
		homepage: homepage,
		tasks: tasks
	},
	mounted: function() {
		window.bus.$on('logged-in', this.login);
		window.bus.$on('signed-up', this.login);
	},
	created: function() {
		this.checkSavedAuth();
		this.returnToPrevious();
	},
	methods: {
		toggleNav: function() {
			this.isNavOpen = !this.isNavOpen;
		},
		loadPage: function(page) {
			this.isNavOpen = false;
			sessionStorage.setItem('previousPage', this.currentView);
			this.currentView = page;
			sessionStorage.setItem('currentPage', page);
		},
		loginModal: function() {
			window.bus.$emit('login-modal');
		},
		signupModal: function() {
			window.bus.$emit('signup-modal');
		},
		login: function(user) {
			this.isLoggedIn = true;
			this.isNavOpen = false;
			this.user = user;
			sessionStorage.setItem('pro', JSON.stringify(user));
			this.loadPage('tasks');
		},
		logout: function() {
			this.isLoggedIn = false;
			this.isNavOpen = false;
			this.user = {};
			sessionStorage.removeItem('pro');
			this.loadPage('homepage');
		},
		checkSavedAuth: function() {
			pre = JSON.parse(sessionStorage.getItem('pro'));
			if (pre && pre.username !== "" && pre.token !== "") {
				this.isLoggedIn = true;
				this.user = pre;
			}
		},
		returnToPrevious: function() {
			if ( ! this.isLoggedIn) { return; }
			last = sessionStorage.getItem('currentPage');
			if (last !== undefined && last !== null && last !== "") {
				this.currentView = last;
			}
		}
	}
});
