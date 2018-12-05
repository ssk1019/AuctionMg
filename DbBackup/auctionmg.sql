/*
Navicat MySQL Data Transfer

Source Server         : AuctionMg
Source Server Version : 50710
Source Host           : localhost:3306
Source Database       : auctionmg

Target Server Type    : MYSQL
Target Server Version : 50710
File Encoding         : 65001

Date: 2018-12-04 19:25:45
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `productinfo`
-- ----------------------------
DROP TABLE IF EXISTS `productinfo`;
CREATE TABLE `productinfo` (
  `ProductId` varchar(20) NOT NULL DEFAULT '' COMMENT '產品ID',
  `Nickname` varchar(255) NOT NULL DEFAULT '',
  `CreateTime` datetime NOT NULL,
  `NameCN` varchar(255) NOT NULL,
  `NameTW` varchar(255) NOT NULL,
  `Weight` float(8,0) NOT NULL,
  `Content` varchar(255) NOT NULL,
  PRIMARY KEY (`ProductId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of productinfo
-- ----------------------------

-- ----------------------------
-- Table structure for `shopitemlist`
-- ----------------------------
DROP TABLE IF EXISTS `shopitemlist`;
CREATE TABLE `shopitemlist` (
  `ShopId` varchar(255) NOT NULL,
  `UpdateTime` datetime NOT NULL,
  `ItemIdList` varchar(255) NOT NULL,
  `ItemIdCnt` int(10) unsigned NOT NULL,
  PRIMARY KEY (`ShopId`),
  KEY `ShopId` (`ShopId`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of shopitemlist
-- ----------------------------
