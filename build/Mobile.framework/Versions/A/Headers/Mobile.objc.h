// Objective-C API for talking to github.com/spolabs/wallet-api/src/api/mobile Go package.
//   gobind -lang=objc github.com/spolabs/wallet-api/src/api/mobile
//
// File is generated by gobind. Do not edit.

#ifndef __Mobile_H__
#define __Mobile_H__

@import Foundation;
#include "Universe.objc.h"


@protocol MobileCoiner;
@class MobileCoiner;

@protocol MobileCoiner <NSObject>
- (NSString*)broadcastTx:(NSString*)rawtx error:(NSError**)error;
// skipped method Coiner.CreateRawTx with unsupported parameter or return types

- (NSString*)getBalance:(NSString*)addrs error:(NSError**)error;
- (NSString*)getNodeAddr;
- (NSString*)getOutputByID:(NSString*)outid error:(NSError**)error;
- (NSString*)getTransactionByID:(NSString*)txid error:(NSError**)error;
- (BOOL)isTransactionConfirmed:(NSString*)txid ret0_:(BOOL*)ret0_ error:(NSError**)error;
- (NSString*)name;
- (NSString*)send:(NSString*)walletID toAddr:(NSString*)toAddr amount:(NSString*)amount error:(NSError**)error;
- (BOOL)validateAddr:(NSString*)addr error:(NSError**)error;
@end

/**
 * GetAddresses return all addresses in the wallet.
returns {"addresses":["jvzYqvdZs17i67cxZ5R8zGE4446JGPVYyz","FNhfaxwWgDVfuXdn2kUoMkxpDFGvqoSPzq","5spraVxAAkFC9j1cpMEdMu7CoV3iHRG7pG"]}
 */
FOUNDATION_EXPORT NSString* MobileGetAddresses(NSString* walletID, NSError** error);

/**
 * GetBalance return balance of a specific address.
returns {"balance":"70.000000"}
 */
FOUNDATION_EXPORT NSString* MobileGetBalance(NSString* coinType, NSString* address, NSError** error);

/**
 * GetKeyPairOfAddr get pubkey and seckey pair of address in specific wallet.
 */
FOUNDATION_EXPORT NSString* MobileGetKeyPairOfAddr(NSString* walletID, NSString* addr, NSError** error);

/**
 * GetSeed returun wallet seed
 */
FOUNDATION_EXPORT NSString* MobileGetSeed(NSString* walletID, NSError** error);

/**
 * GetTransactionByID gets transaction verbose info by id
 */
FOUNDATION_EXPORT NSString* MobileGetTransactionByID(NSString* coinType, NSString* txid, NSError** error);

/**
 * GetWalletBalance return balance of wallet.
 */
FOUNDATION_EXPORT NSString* MobileGetWalletBalance(NSString* coinType, NSString* wltID, NSError** error);

/**
 * Init initialize wallet dir and coin manager.
 */
FOUNDATION_EXPORT void MobileInit(NSString* walletDir);

/**
 * IsContain wallet contains address (format "a1,a2,a3") or not
 */
FOUNDATION_EXPORT BOOL MobileIsContain(NSString* walletID, NSString* addrs, BOOL* ret0_, NSError** error);

/**
 * IsExist wallet exists or not
 */
FOUNDATION_EXPORT BOOL MobileIsExist(NSString* walletID);

/**
 * IsTransactionConfirmed gets transaction verbose info by id
 */
FOUNDATION_EXPORT BOOL MobileIsTransactionConfirmed(NSString* coinType, NSString* txid, BOOL* ret0_, NSError** error);

/**
 * NewAddress generate address in specific wallet.
 */
FOUNDATION_EXPORT NSString* MobileNewAddress(NSString* walletID, long num, NSError** error);

/**
 * NewSeed generates mnemonic seed
 */
FOUNDATION_EXPORT NSString* MobileNewSeed(void);

/**
 * NewWallet create a new wallet base on the wallet type and seed
 */
FOUNDATION_EXPORT NSString* MobileNewWallet(NSString* coinType, NSString* lable, NSString* seed, NSError** error);

/**
 * RegisterNewCoin register a new coin to wallet
the server address is consisted of ip and port, eg: 127.0.0.1:6420
 */
FOUNDATION_EXPORT BOOL MobileRegisterNewCoin(NSString* coinType, NSString* serverAddr, NSError** error);

/**
 * Remove delete wallet.
 */
FOUNDATION_EXPORT BOOL MobileRemove(NSString* walletID, NSError** error);

/**
 * Send send coins, support bitcoin and all coins in skycoin ledger
 */
FOUNDATION_EXPORT NSString* MobileSend(NSString* coinType, NSString* wid, NSString* toAddr, NSString* amount, NSError** error);

/**
 * ValidateAddress validate the address
 */
FOUNDATION_EXPORT BOOL MobileValidateAddress(NSString* coinType, NSString* addr, BOOL* ret0_, NSError** error);

@class MobileCoiner;

/**
 * Coiner coin client interface
 */
@interface MobileCoiner : NSObject <goSeqRefInterface, MobileCoiner> {
}
@property(strong, readonly) id _ref;

- (instancetype)initWithRef:(id)ref;
- (NSString*)broadcastTx:(NSString*)rawtx error:(NSError**)error;
// skipped method Coiner.CreateRawTx with unsupported parameter or return types

- (NSString*)getBalance:(NSString*)addrs error:(NSError**)error;
- (NSString*)getNodeAddr;
- (NSString*)getOutputByID:(NSString*)outid error:(NSError**)error;
- (NSString*)getTransactionByID:(NSString*)txid error:(NSError**)error;
- (BOOL)isTransactionConfirmed:(NSString*)txid ret0_:(BOOL*)ret0_ error:(NSError**)error;
- (NSString*)name;
- (NSString*)send:(NSString*)walletID toAddr:(NSString*)toAddr amount:(NSString*)amount error:(NSError**)error;
- (BOOL)validateAddr:(NSString*)addr error:(NSError**)error;
@end

#endif
