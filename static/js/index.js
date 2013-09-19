(function() {

var Welcome = Backbone.View.extend({
  el: '.welcome',
  events: {
    'submit form': 'create'
  },
  create: function(e) {
    e.preventDefault();
    var user = this.$('.user').val();
    if (user.length > 0) {
      this.$el.hide();
      this.trigger('create', user);
    }
  }
});

var Join = Backbone.View.extend({
  el: '.join',
  events: {
    'submit form': 'join'
  },
  join: function(e) {
    e.preventDefault();
    var user = this.$('.user').val();
    if (user.length > 0) {
      this.$el.hide();
      this.trigger('join', user, this.model.get('name'));
    }
  }
});

var Room = Backbone.View.extend({
  el: '.room',
  events: {
    'submit form': 'send'
  },
  initialize: function(options) {
    _.bindAll(this, 'handleJoin', 'handleMsg');
    this.conn = options.conn;
    this.setupConnListeners();
    this.setupModelListeners();
  },
  setupConnListeners: function() {
    this.listenTo(this.conn, 'join_room', this.handleJoin);
    this.listenTo(this.conn, 'leave_room', this.handleLeave);
    this.listenTo(this.conn, 'msg', this.handleMsg);
  },
  handleJoin: function(data) {
    this.notice({
      text: "Andito na si " + data.user + "!"
    });
    this.$('.msg').focus();
  },
  handleLeave: function(data) {
    this.notice({
      text: "Umalis na si " + data.user + "!"
    });
  },
  handleMsg: function(data) {
    this.message(data);
  },
  setupModelListeners: function () {
    this.listenTo(this.model, 'change:user', this.updateUser);
  },
  updateUser: function() {
    this.$('h2').text(this.model.get('user'));
  },
  notice: function(msg) {
    this.$('.log ul').append(
      $('<li>').addClass('notice').text(msg.text)
    );
    this.scrollBot();
  },
  message: function(msg) {
    this.$('.log ul').append(
      $('<li>')
        .addClass('message')
        .text(msg.text)
        .prepend(
          $('<strong>').text(msg.user)
        )
    );
    this.scrollBot();
  },
  scrollBot: function() {
    var $list = this.$('.log ul');
    $list.get(0).scrollTop = $list.attr('scrollHeight');
  },
  send: function(e) {
    e.preventDefault();
    var $msg = this.$('.msg');
    var text = $msg.val();
    if (text.length > 0) {
      this.conn.emit('msg', {
        room: this.model.get('name'),
        text: text
      });
      $msg.val('');
    }
  }
});

var App = Backbone.View.extend({
  initialize: function(options) {
    _.bindAll(this, 'handleJoin');
    this.connect();
    this.router = options.router;
    this.router.app = this;
    this.room = new Backbone.Model();
    this.views = {
      welcome: new Welcome(),
      join: new Join({model: this.room}),
      room: new Room({model: this.room, conn: this.conn})
    };
    this.setupViewListeners();
  },
  connect: function() {
    var o = location.origin.split(':');
    o.shift();
    o.unshift('ws');
    o.join(':') + '/ws'
    this.conn = new golem.Connection(o.join(':') + '/ws');
    this.setupConnListeners();
  },
  setupConnListeners: function() {
    this.listenTo(this.conn, 'join', this.handleJoin);
  },
  setupViewListeners: function() {
    this.listenTo(this.views.welcome, 'create', this.join);
    this.listenTo(this.views.join, 'join', this.join);
  },
  join: function(user, room) {
    this.conn.emit('join', {user: user, room: room});
  },
  handleJoin: function(data) {
    this.room.set({
      name: data.room,
      user: data.user
    });
    this.views.room.$el.show();
    this.router.navigate(data.room);
  }
})

var Router = Backbone.Router.extend({
  routes: {
    '': 'welcome',
    ':room': 'join'
  },
  welcome: function() {
    this.app.views.welcome.$el.show().find('.user').focus();
  },
  join: function(room) {
    this.app.room.set('name', room);
    this.app.views.join.$el.show().find('.user').focus();
  }
});


$(function() {
  var app = new App({router: new Router()});
  Backbone.history.start({pushState: true});
});

})();
