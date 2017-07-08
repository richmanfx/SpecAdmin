-- MySQL dump 10.13  Distrib 5.7.17, for Win64 (x86_64)
--
-- Host: localhost    Database: specadmin
-- ------------------------------------------------------
-- Server version	5.7.17-log
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `configuration`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `configuration` (
  `id` int(11) NOT NULL,
  `property_name` varchar(100) DEFAULT NULL,
  `property_value` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `configuration_property_name_uindex` (`property_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tests_groups`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tests_groups` (
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Группы тестов';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tests_scripts`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tests_scripts` (
  `id` int(11) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `serial_number` int(11) DEFAULT NULL,
  `name_suite` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `tests_scripts_tests_suits__fk` (`name_suite`),
  CONSTRAINT `tests_scripts_tests_suits__fk` FOREIGN KEY (`name_suite`) REFERENCES `tests_suits` (`name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tests_steps`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tests_steps` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `serial_number` int(11) NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `expected_result` varchar(500) DEFAULT NULL,
  `script_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `tests_steps_tests_scripts__fk` (`script_id`),
  CONSTRAINT `tests_steps_tests_scripts__fk` FOREIGN KEY (`script_id`) REFERENCES `tests_scripts` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Шаги тестов';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tests_suits`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tests_suits` (
  `name` varchar(255) NOT NULL,
  `serial_number` int(11) NOT NULL,
  `description` varchar(500) NOT NULL,
  `name_group` varchar(255) NOT NULL,
  PRIMARY KEY (`name`),
  KEY `tests_suits_tests_groups__fk` (`name_group`),
  CONSTRAINT `tests_suits_tests_groups__fk` FOREIGN KEY (`name_group`) REFERENCES `tests_groups` (`name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Сюиты тестов';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'specadmin'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-07-06 18:36:22
