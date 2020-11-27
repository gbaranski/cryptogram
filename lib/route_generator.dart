import 'package:cryptogram/views/error_screen.dart';

import 'views/index.dart';
import 'views/account/index.dart';
import 'views/account/restore.dart';
import 'views/chat/index.dart';
import 'package:flutter/material.dart';

class RouteGenerator {
  static Route<dynamic> generateRoute(RouteSettings settings) {
    return MaterialPageRoute(builder: (context) {
      switch (settings.name) {
        case IndexView.route:
          return IndexView();
        case ChatView.route:
          return ChatView();
        case RestoreAccountView.route:
          return RestoreAccountView();
        case AccountView.route:
          return AccountView();
        default:
          return ErrorScreen();
      }
    });
  }
}
