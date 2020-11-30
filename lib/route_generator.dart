import 'package:cryptogram/models/account.dart';
import 'package:cryptogram/views/account/accounts_list.dart';
import 'package:cryptogram/views/account/create.dart';
import 'package:cryptogram/views/error_screen.dart';
import 'package:cryptogram/views/index.dart';

import 'views/account/restore.dart';
import 'package:flutter/material.dart';

class RouteGenerator {
  static Route<dynamic> generateRoute(RouteSettings settings) {
    return MaterialPageRoute(
        builder: (context) {
          switch (settings.name) {
            case IndexView.route:
              final Account account = settings.arguments;
              if (account == null) return AccountsList();
              return IndexView(account);
            case AccountsList.route:
              return AccountsList();
            case CreateAccountView.route:
              return CreateAccountView();
            case RestoreAccountView.route:
              return RestoreAccountView();
            default:
              return ErrorScreen();
          }
        },
        settings: settings);
  }
}
