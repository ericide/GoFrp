//
//  JLTouchUnlock.h
//  touchUnlock
//
//  Created by junlei on 15/7/6.
//  Copyright (c) 2015年 HA. All rights reserved.
//

#import <UIKit/UIKit.h>
@class JLTouchUnlock;

@protocol JLTouchUnlockDelegate
@optional
- (void)secureTextWithTarget:(JLTouchUnlock *)target text:(NSString *)str;
- (void)startInputText:(JLTouchUnlock *)target;
@end

@interface JLTouchUnlock : UIView
@property (nonatomic,weak) id<JLTouchUnlockDelegate> delegate;
@property (nonatomic,assign) BOOL lineEnable;
@end
