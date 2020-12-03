import 'package:cryptogram/models/account.dart';
import 'package:cryptogram/services/database.dart';
import 'package:cryptogram/views/account/create.dart';
import 'package:cryptogram/views/index.dart';
import 'package:flutter/material.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

enum AccountAction {
  delete,
  edit,
}

class AccountsList extends StatefulWidget {
  static const route = '/accounts';

  @override
  _AccountsListState createState() => _AccountsListState();
}

class _AccountsListState extends State<AccountsList> {
  List<Account> _accounts; // to be implemented

  Future<void> loadAccounts() async {
    final accounts = await DatabaseService.getAccounts();
    print("Loaded accounts, count: ${accounts.length}");
    setState(() {
      _accounts = accounts;
    });
  }

  Future<void> onAccountAction(AccountAction action, int index) async {
    if (action == AccountAction.delete) {
      final targetAccount = _accounts[index];
      if (targetAccount == null)
        throw new Exception(
            'Exception while removing account: Cannot find account');
      await DatabaseService.deleteAccount(targetAccount);
      await loadAccounts();
    }
  }

  void onSelected(BuildContext context, int index) {
    final selectedAccount = _accounts[index];
    if (selectedAccount == null)
      throw new Exception(
          "Error when navigating to chat screen: Couldn't find proper accounnt");
    Navigator.pushReplacementNamed(
      context,
      IndexView.route,
      arguments: selectedAccount,
    );
  }

  @override
  void initState() {
    super.initState();
    loadAccounts();
  }

  Widget accountsList(BuildContext context) => SliverList(
          delegate: SliverChildBuilderDelegate((context, i) {
        final account = _accounts[i];
        return ListTile(
          onTap: () => onSelected(context, i),
          leading: Icon(MdiIcons.accountKey),
          trailing: PopupMenuButton<AccountAction>(
            onSelected: (action) => onAccountAction(action, i),
            icon: Icon(Icons.more_vert),
            itemBuilder: (context) => [
              PopupMenuItem(
                  value: AccountAction.delete,
                  child: ListTile(
                    visualDensity: VisualDensity.compact,
                    contentPadding: EdgeInsets.zero,
                    leading: Icon(Icons.delete),
                    title: Text("Delete"),
                  )),
              PopupMenuItem(
                  value: AccountAction.edit,
                  enabled: false,
                  child: ListTile(
                    visualDensity: VisualDensity.compact,
                    contentPadding: EdgeInsets.zero,
                    leading: Icon(Icons.edit),
                    title: Text(
                      "Edit",
                      style: TextStyle(color: Colors.white60),
                    ),
                  )),
            ],
          ),
          title: Text(
            account.accountID,
            overflow: TextOverflow.ellipsis,
          ),
          subtitle: Text(account.customName),
        );
      }, childCount: _accounts.length));

  Future<void> navigateToCreateAccount() async {
    final result = await Navigator.pushNamed(context, CreateAccountView.route)
        as CreateAccountResult;
    if (result != null && result.ok) loadAccounts();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: FloatingActionButton(
        onPressed: navigateToCreateAccount,
        child: Icon(MdiIcons.accountPlus),
      ),
      body: RefreshIndicator(
        onRefresh: loadAccounts,
        child: CustomScrollView(
          physics:
              AlwaysScrollableScrollPhysics().applyTo(BouncingScrollPhysics()),
          slivers: [
            SliverAppBar(
              title: Text("Pick account"),
            ),
            if (_accounts != null)
              accountsList(context)
            else
              (SliverToBoxAdapter(
                child: CircularProgressIndicator(),
              ))
          ],
        ),
      ),
    );
  }
}