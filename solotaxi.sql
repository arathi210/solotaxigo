-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Nov 25, 2022 at 10:05 AM
-- Server version: 10.4.25-MariaDB
-- PHP Version: 7.4.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `solotaxi`
--

-- --------------------------------------------------------

--
-- Table structure for table `custom_fare`
--

CREATE TABLE `custom_fare` (
  `id` int(11) NOT NULL,
  `fare_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `base_fare` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `min_fare` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `cost` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `custom_fare`
--

INSERT INTO `custom_fare` (`id`, `fare_name`, `base_fare`, `min_fare`, `cost`, `user_id`) VALUES
(1, 'Day Ride', '24', '18', '20', 5),
(2, 'Day Ride1', '24', '18', '20', 1),
(3, 'Day Ride3', '24', '18', '20', 1),
(4, 'Day Ride4', '24', '18', '20', 1),
(5, 'Day Ride6', '34', '38', '30', 5),
(6, 'Day Ride666', '111', '111', '111', 1),
(7, 'Day Ride44', '1', '1', '1', 1),
(8, 'Day Ride23232', '11', '11', '11', 5),
(9, 'Day Ride222', '112', '112', '22', 5);

-- --------------------------------------------------------

--
-- Table structure for table `ride_history`
--

CREATE TABLE `ride_history` (
  `id` int(11) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `from_lat` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `from_lon` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `from_address` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `from_date` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `from_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `base_fare` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cost` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `to_lat` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `to_lon` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `to_address` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `to_date` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `to_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `distance` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `waiting` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `total` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `duration` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount_status` tinyint(1) DEFAULT 0
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `ride_history`
--

INSERT INTO `ride_history` (`id`, `user_id`, `from_lat`, `from_lon`, `from_address`, `from_date`, `from_time`, `base_fare`, `cost`, `to_lat`, `to_lon`, `to_address`, `to_date`, `to_time`, `distance`, `waiting`, `total`, `duration`, `amount_status`) VALUES
(19, 1, '15.9072815', '74.5192699', 'WG49 VM6, Kangrali, Belagavi, Karnataka 590010, India', '11/24/2022', '18:14', '24', '20', '16.164269', '74.819126', '5R79+PJ9, Gokak, Karnataka 591307, India', '11-24-2022', '18:14', '42.98', '12', '872', '00 Hour 00 Min ', 0),
(20, 1, '15.9072815', '74.5192699', 'WG49 VM6, Kangrali, Belagavi, Karnataka 590010, India', '11/24/2022', '18:17', '24', '20', '16.164269', '74.819126', '5R79+PJ9, Gokak, Karnataka 591307, India', '11-24-2022', '18:17', '42.98', '12', '872', '00 Hour 00 Min ', 1);

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mobile` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `name`, `mobile`, `email`, `password`) VALUES
(1, 'Admin', '9999999999', 'admin@yopmail.com', '$2a$10$t5JlSMEArrnY29tKYI2w1OwkFJyiUIiFKpmfn6cJwDoZQa/5Q1Twe'),
(5, 'Admin2', '9900798019', 'admin2@yopmail.com', '$2a$10$Vxr.n449ViqA5U.avs.83e11C0Wmnz0xx6pNOz9vKdtCeFgddypwS');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `custom_fare`
--
ALTER TABLE `custom_fare`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `ride_history`
--
ALTER TABLE `ride_history`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `users_email_unique` (`email`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `custom_fare`
--
ALTER TABLE `custom_fare`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT for table `ride_history`
--
ALTER TABLE `ride_history`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `custom_fare`
--
ALTER TABLE `custom_fare`
  ADD CONSTRAINT `custom_fare_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
