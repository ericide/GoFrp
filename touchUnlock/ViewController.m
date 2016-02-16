//
//  ViewController.m
//  touchUnlock
//
//  Created by junlei on 15/7/6.
//  Copyright (c) 2015年 HA. All rights reserved.
//

#import "ViewController.h"
#import "JLTouchUnlock.h"

@interface ViewController ()<JLTouchUnlockDelegate>

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view, typically from a nib.
    JLTouchUnlock * touchUnlock = [[JLTouchUnlock alloc]initWithFrame:[[UIScreen mainScreen] bounds]];
    touchUnlock.delegate = self;
    touchUnlock.lineEnable = YES;
    [self.view addSubview:touchUnlock];

}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}
- (void)secureTextWithTarget:(JLTouchUnlock *)target text:(NSString *)str
{
    NSLog(@"%@",str);
}
- (void)startInputText:(JLTouchUnlock *)target
{
    NSLog(@"开始了");
}
@end
