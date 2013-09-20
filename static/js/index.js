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
    _.bindAll(
      this,
      'handleJoin',
      'handleLeave',
      'handleMsg',
      'handleUsers',
      'handleClose'
    );
    this.conn = options.conn;
    this.setupConnListeners();
    this.setupModelListeners();
  },
  setupConnListeners: function() {
    this.listenTo(this.conn, 'join_room', this.handleJoin);
    this.listenTo(this.conn, 'leave_room', this.handleLeave);
    this.listenTo(this.conn, 'msg', this.handleMsg);
    this.listenTo(this.conn, 'users', this.handleUsers);
    this.listenTo(this.conn, 'close', this.handleClose);
  },
  handleJoin: function(data) {
    this.notice({
      icon: 'icon-login',
      text: data.user
    });
    this.$('.msg').focus();
  },
  handleLeave: function(data) {
    this.notice({
      icon: 'icon-logout',
      text: data.user
    });
  },
  handleMsg: function(data) {
    this.message(data);
  },
  handleUsers: function(data) {
    var $list = this.$('.users ul').empty();
    _.each(data.users, function(user) {
      $list.append(
        $('<li>')
          .text(user)
          .prepend(
            $('<i>').addClass('icon-smiley')
          )
      );
    });
  },
  handleClose: function(data) {
    this.notice({
      icon: 'icon-warning',
      text: 'Got disconnected.'
    });
  },
  setupModelListeners: function () {
    this.listenTo(this.model, 'change', this.updateTitle);
  },
  updateTitle: function() {
    this.$('h2 span.name').text(this.model.get('name'));
    this.$('h2 span.user').text(this.model.get('user'));
  },
  notice: function(msg) {
    this.$('.log ul').append(
      $('<li>')
        .addClass('notice')
        .text(msg.text)
        .prepend(
          $('<i>').addClass(msg.icon)
        )
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
  el: '.app',
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
    this.setupModelListeners();
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
  setupModelListeners: function() {
    this.listenTo(this.room, 'change:name', this.updateTitle);
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
  },
  updateTitle: function() {
    document.title = this.room.get('name');
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
