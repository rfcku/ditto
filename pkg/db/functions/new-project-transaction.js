exports = async function (changeEvent) {
  const doc = changeEvent.fullDocument;
  console.log('test', doc);

  try {
    if (changeEvent.operationType === 'insert') {
      const collection = context.services
        .get('mongodb-atlas')
        .db('spark')
        .collection('wallets');

      const wallet = await collection.updateOne(
        { _id: doc.author_id },
        {
          $set: {
            balance: doc.balance - 10,
            updatedAt: new Date(),
          },
        }
      );

      const collection2 = context.services
        .get('mongodb-atlas')
        .db('spark')
        .collection('transactions');

      await collection2.insertOne({
        balance_prev: wallet.balance,
        amount: 10,
        type: 'post',
        user_id: wallet.author_id,
        post_id: doc._id,
        createdAt: new Date(),
      });

      console.log('This was called');
    }
  } catch (err) {
    console.log('error performing mongodb write: ', err.message);
  }
};
