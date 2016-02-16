//
//  JLTouchUnlock.m
//  touchUnlock
//
//  Created by junlei on 15/7/6.
//  Copyright (c) 2015年 HA. All rights reserved.
//

#define PICY 74
#define LOCATIONY ([UIScreen mainScreen].bounds.size.height*0.3)

#import "JLTouchUnlock.h"
@interface JLTouchUnlock()
@property (nonatomic,strong) NSMutableArray * touchElements;
@property (nonatomic,assign) CGPoint  currentPoint;
@end
@implementation JLTouchUnlock

-(id)initWithCoder:(NSCoder *)aDecoder
{
    self = [super initWithCoder:aDecoder];
    if (self != nil) {
        [self setup];
    }
    return self;
}
- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self != nil) {
        [self setup];
    }
    return self;
}
- (void)setup
{
    self.lineEnable = YES;
    for (int num = 0; num < 9; num ++) {
        UIButton * but = [[UIButton alloc]initWithFrame:CGRectMake(0, 0, 74, 74)];
        but.tag = 1024+num;
        [but setImage:[UIImage imageNamed:@"gesture_node_normal.png"] forState:UIControlStateNormal];
        [but setImage:[UIImage imageNamed:@"gesture_node_highlighted.png"] forState:UIControlStateSelected];
        but.userInteractionEnabled = NO;
        [self addSubview:but];
    }
}
- (void)layoutSubviews
{
    float width = [UIScreen mainScreen].bounds.size.width;
    float widthW = width/3.5;
    float heightH = widthW;
    float offsetX = width/2 - widthW ;
    float offsetY = PICY/2 + LOCATIONY;
    for (int num = 0; num < 9; num ++) {
        UIButton * but = (UIButton *)[self viewWithTag:(1024+num)];
        but.center = CGPointMake(widthW*(num%3)+offsetX, (num/3)*heightH+offsetY);
    }
    
}
- (NSMutableArray *)touchElements
{
    if (_touchElements == nil) {
        _touchElements = [NSMutableArray array];
    }
    return _touchElements;
}
- (CGPoint )pointWithTouch:(NSSet *)touches
{
    UITouch * touch = [touches anyObject];
    return [touch locationInView:touch.view];
}
- (UIButton *)buttonWithPoint:(CGPoint)point
{
    UIButton * selectBut = nil;
    for (UIButton * but in self.subviews) {
        if (CGRectContainsPoint(but.frame, point) ) {
            selectBut = but;
            break;
        }
    }
    return selectBut;
}
- (void)touchesBegan:(NSSet *)touches withEvent:(UIEvent *)event
{
    
    CGPoint touchPoint = [self pointWithTouch:touches];
    self.currentPoint = touchPoint;
    UIButton * touchBut = [self buttonWithPoint:touchPoint];
    
    if (touchBut != nil && !touchBut.selected) {
        [self.touchElements addObject:touchBut];
        touchBut.selected = YES;
    }
    
    [self setNeedsDisplay];
    
    [self.delegate startInputText:self];
}
- (void)touchesMoved:(NSSet *)touches withEvent:(UIEvent *)event
{
    CGPoint touchPoint = [self pointWithTouch:touches];
    self.currentPoint = touchPoint;
    UIButton * touchBut = [self buttonWithPoint:touchPoint];

    if (touchBut != nil && !touchBut.selected) {
        [self.touchElements addObject:touchBut];
        touchBut.selected = YES;
    }
    
    [self setNeedsDisplay];
}
- (void)touchesEnded:(NSSet *)touches withEvent:(UIEvent *)event
{
    NSMutableString * str = [NSMutableString string];
    for (int num = 0 ;num <self.touchElements.count ; num++) {
        UIButton * but = self.touchElements[num];
        [str appendFormat:@"%ld",(but.tag-1024)];
        but.selected = NO;
    }
    
    [self.touchElements removeAllObjects];
    [self setNeedsDisplay];
    
    [self.delegate secureTextWithTarget:self text:str];
}
- (BOOL)isOpaque
{
    return NO;
}
// Only override drawRect: if you perform custom drawing.
// An empty implementation adversely affects performance during animation.
- (void)drawRect:(CGRect)rect {
    // Drawing code
    if (!self.lineEnable) {
        return;
    }
    UIBezierPath * path = [UIBezierPath bezierPath];
    
    for (int num = 0; num < self.touchElements.count; num++) {
        UIView * view = self.touchElements[num];
        if (num == 0) {
            [path moveToPoint:view.center];
        } else {
            [path addLineToPoint:view.center];
        }
        
    }
    
    if (self.touchElements.count > 0 ) {
        [path addLineToPoint:self.currentPoint];
    }
    path.lineWidth = 8;
    path.lineJoinStyle = kCGLineJoinRound;
    [[UIColor colorWithRed:32/255.0 green:210/255.0 blue:254/255.0 alpha:0.5] set];
    [path stroke];
    
}


@end
