<?php

session_start();

if (!isset($_SESSION['logged_admin'])) {
    header('Location: index.php');
    exit();
}
require_once "../connect.php";
$_SESSION['conn'] = @new mysqli($host, $db_user, $db_password, $db_name);
?>
<!DOCTYPE html>
<html lang="pl">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Zbiórki</title>
    <link rel="stylesheet" href="../general.css">
    <link rel="stylesheet" href="zbiorki.css">
</head>

<body>
    <main>
        <div class="container">
            <button class="showbutton" onclick="show_nav()">
                <ion-icon name="menu-outline" class="hide620"></ion-icon>
            </button>
            <nav id="main_nav">

                <ul>
                    <li><a href="../home.php">
                            <ion-icon name="home-outline"></ion-icon><span class="menu">Dom</span>
                        </a></li>

                    <li>
                        <a href="../frekwencja.php#week">
                            <ion-icon name="clipboard-outline"></ion-icon><span class="menu">Frekwencja</span>
                        </a>
                    </li>
                    <li><a href="../ranking.php">
                            <ion-icon name="trophy-outline"></ion-icon><span class="menu">Ranking</span>
                        </a>
                    </li>
                    <li class="checked">
                        <a href="#">
                            <ion-icon name="reader-outline"></ion-icon><span class="menu">Zbiórki</span>
                        </a>
                    </li>
                    <li>
                        <a href="http://sluzbaliturgiczna-zgloszenia.epizy.com/help.php">
                            <ion-icon name="information-circle-outline"></ion-icon><span class="menu">Pomoc</span>
                        </a>
                    </li>

                    <li>
                        <a href="logout.php">
                            <ion-icon name="log-out-outline"></ion-icon><span class="menu">Wyloguj się</span>
                        </a>
                    </li>
                    <li class="hide1030">
                        <button onclick="hide_nav()">
                            <a>
                                <ion-icon name="chevron-forward-circle-outline" class="zwin"></ion-icon><span class="menu">Zwiń</span>
                            </a>
                        </button>
                    </li>
                </ul>


            </nav>
            <article class="II_blocks zbiorki_article">
                <header>
                <ul>
                        <li class="top_link"><a href="addObecnosc.php">Dodaj obecność</a></li>
                    </ul>
                    <ul>
                        <li class="top_link"><a href="usprawiedliwienia.php">Usprawiedliwienia</a></li>
                    </ul>
                    <ul>
                        <li class="top_link checked"><a href="">Zbiórki</a></li>
                    </ul>
                    <ul>
                        <li><a href="zbiorki_lektor.php">Lektorzy</a></li>
                        <li class="checked"><a href="zbiorki_ministrant.php">Ministranci</a></li>
                    </ul>
                </header>
                <form action="verify.php" method="POST" id="lektorzy">
                    <ul>
                        <?php
                        $result = $_SESSION['conn']->query("SELECT Imie,Nazwisko,karta FROM `osoby` WHERE stopien LIKE 'M%';");
                        if ($result->num_rows > 0) {
                            while ($row = $result->fetch_assoc()) {
                                echo '
                                <li>
                                    <p>' . $row['Imie'] . ' ' . $row['Nazwisko'] . '</p>
                                    <label>
                                        <span>
                                            <input type="radio" name="' . $row['karta'] . '" id="" value="present" required class="present"><ion-icon name="checkmark-outline"></ion-icon>
                                        </span>
                                        <span>
                                            <input type="radio" name="' . $row['karta'] . '" id="" value="justified" required class="justified"><ion-icon name="chatbox-ellipses-outline"></ion-icon>
                                        </span>
                                        <span>
                                            <input type="radio" name="' . $row['karta'] . '" id="" value="absent" required class="absent"><ion-icon name="close-outline"></ion-icon>
                                        </span>
                                    </label>
                                </li>
                                ';
                            }
                        }
                        ?>

                    </ul>
                    <input type="hidden" name="stopien" value="M">
                    <input type="submit" value="Zapisz">
                </form>























            </article>

        </div>
    </main>
    <script src="../button_script.js"></script>
    <script type="module" src="https://unpkg.com/ionicons@5.5.2/dist/ionicons/ionicons.esm.js"></script>
    <script nomodule src="https://unpkg.com/ionicons@5.5.2/dist/ionicons/ionicons.js"></script>
</body>

</html>