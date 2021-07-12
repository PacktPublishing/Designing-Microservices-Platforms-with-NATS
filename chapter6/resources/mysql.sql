CREATE TABLE `patient_details` (
`id` int(15) unsigned NOT NULL,
`full_name` varchar(100) NOT NULL,
`address` varchar(255),
`sex` varchar(10),
`phone` int(15) unsigned,
`remarks` varchar(255),
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `patient_registrations` (
`id` int(15) unsigned NOT NULL,
`token` int(10) unsigned NOT NULL,
PRIMARY KEY (`token`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `inspection_reports` (
`id` int(15) unsigned NOT NULL,
`medication` varchar(255),
`tests` varchar(255),
`notes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `inspection_details` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50) NOT NULL,
`observations` varchar(255),
`medication` varchar(255),
`tests` varchar(255),
`notes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `medication_reports` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`dose` varchar(255),
`notes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `test_reports` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`test_name` varchar(100),
`results` varchar(255),
`status` varchar(50),
`notes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `release_reports` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`next_state` varchar(50),
`post_medication` varchar(255),
`notes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `discharge_details` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`state` varchar(50),
`post_medication` varchar(255),
`notes` varchar(255),
`next_visit` varchar(50)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
