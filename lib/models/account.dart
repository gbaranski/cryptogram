class Account {
  final String accountID;
  final String customName;

  Map<String, dynamic> toMap() {
    return {
      'accountID': accountID,
      'customName': customName,
    };
  }

  Account({this.accountID, this.customName});
}
