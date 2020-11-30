import 'package:cryptogram/models/account.dart';
import 'package:cryptogram/services/blockchain.dart';
import 'package:cryptogram/views/account/accounts_list.dart';
import 'package:flutter/material.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

class AccountView extends StatelessWidget {
  final Account account;

  AccountView(this.account);

  Future<void> executeAction(BuildContext context, dynamic Function() action,
      [String onSuccess]) async {
    ScaffoldMessenger.of(context)
        .hideCurrentSnackBar(reason: SnackBarClosedReason.timeout);

    try {
      await action();
      if (onSuccess != null)
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          content: Text(onSuccess),
        ));
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(e.toString()),
      ));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          SliverAppBar(),
          SliverList(
            delegate: SliverChildListDelegate([
              ListTile(
                onTap: () {},
                isThreeLine: true,
                leading: Icon(
                  MdiIcons.accountKey,
                ),
                title: Text(account.customName),
                subtitle: Text(account.accountID),
              ),
              ListTile(
                onTap: () => executeAction(
                    context, () => Blockchain.openAccountInfo(account), null),
                leading: Icon(
                  MdiIcons.fileFind,
                ),
                title: Text("View in blockchain explorer"),
              ),
              ListTile(
                onTap: () => executeAction(context, () async {
                  if (!(await Blockchain.fundAccount(account)))
                    throw new Exception('Failed funding account');
                }, "Successfully funded account"),
                leading: Icon(
                  MdiIcons.walletPlus,
                ),
                title: Text("Fund account"),
              ),
              ListTile(
                onTap: () =>
                    Navigator.pushReplacementNamed(context, AccountsList.route),
                leading: Icon(
                  MdiIcons.logout,
                ),
                title: Text("Log out"),
              )
            ]),
          )
        ],
      ),
    );
  }
}
