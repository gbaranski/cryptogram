import 'package:flutter/material.dart';
import 'package:cryptogram/models/account.dart';
import 'chat/index.dart';
import 'account/index.dart';

class NavigationItem {
  IconData iconData;
  String label;
  Widget Function(Account) widget;

  NavigationItem(
      {@required this.iconData, @required this.label, @required this.widget});
}

class IndexView extends StatefulWidget {
  static const route = '/';

  final Account account;
  IndexView(this.account);

  @override
  _IndexViewState createState() => _IndexViewState();
}

class _IndexViewState extends State<IndexView> {
  static final List<NavigationItem> navigationRoutes = [
    NavigationItem(
        iconData: Icons.chat,
        label: 'Chat',
        widget: (Account account) => ChatView(account)),
    NavigationItem(
        iconData: Icons.account_circle,
        label: 'Account',
        widget: (Account account) => AccountView(account))
  ];
  int _selected = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      bottomNavigationBar: BottomNavigationBar(
        items: navigationRoutes
            .map((route) => BottomNavigationBarItem(
                icon: Icon(route.iconData), label: route.label))
            .toList(),
        currentIndex: _selected,
        onTap: (i) => setState(() => _selected = i),
      ),
      body: navigationRoutes[_selected].widget(widget.account),
    );
  }
}
