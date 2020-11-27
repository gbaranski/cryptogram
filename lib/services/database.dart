import 'package:cryptogram/models/account.dart';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';

class DatabaseService {
  static Future<Database> _database;

  static Future<void> init() async {
    _database = openDatabase(join(await getDatabasesPath(), 'database.db'),
        onCreate: (db, version) {
      return db.execute(
          'CREATE TABLE accounts(accountID TEXT PRIMARY KEY, secretSeed TEXT, customName TEXT)');
    }, version: 1);
  }

  static Future<List<Account>> getAccounts() async {
    final List<Map<String, dynamic>> accounts =
        await (await _database).query('accounts');
    return accounts.map((account) => Account.fromMap(account)).toList();
  }

  static Future<void> addAccount(Account account) async {
    (await _database).insert('accounts', account.toMap(),
        conflictAlgorithm: ConflictAlgorithm.fail);
  }

  static Future<void> deleteAccount(Account account) async {
    final int result = await (await _database).delete('accounts',
        where: "accountID = ? AND customName = ? AND secretSeed = ?",
        whereArgs: [account.accountID, account.customName, account.secretSeed]);
    if (result == 0) throw new Exception('Account not found');
  }
}
