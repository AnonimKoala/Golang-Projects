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
    <title>Dodaj obecność</title>
    <link rel="stylesheet" href="../general.css">
    <link rel="stylesheet" href="zbiorki.css">
    <link rel="stylesheet" href="usprwaw.css">
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
                        <a href="zbiorki_lektor.php">
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
                        <li class="top_link"><a href="zbiorki_lektor.php">Powrót</a></li>
                    </ul>


                </header>
                <form action="" method="POST" id="usprawiedliw">
                    <header>
                        <ul>
                            <li class="checked top_link head_uspraw"><a href="">Dodaj obecność</a></li>
                        </ul>
                    </header>
                    <div class="custom-select">

                        <select name="select_person" id="" required>
                            <option value="0">Wybierz osobę</option>
                            <?php
                            $result = $_SESSION['conn']->query('SELECT Imie,Nazwisko,karta, stopien FROM `osoby` WHERE stopien IS NOT NULL AND NOT stopien = "" AND NOT stopien = "A" ORDER BY stopien, Nazwisko ASC');
                            if ($result->num_rows > 0) {
                                while ($row = $result->fetch_assoc()) {
                                    echo '
                                <option value="' . $row['karta'] . '">
                                ' . $row['Imie'] . ' ' . $row['Nazwisko'] . ' &nbsp;&nbsp;&nbsp;[' . $row['stopien'] . ']
                                </option>
                                ';
                                }
                            }
                            ?>

                        </select>

                    </div>
                    <div class="custom-select">
                        <select name="select_time" id="" required>
                            <option value="0">Wybierz godzinę Mszy</option>
                            <option value="6-30">6:30</option>
                            <option value="7-00">7:00</option>
                            <option value="8-00">8:00</option>
                            <option value="9-30">9:30</option>
                            <option value="11-00">11:00</option>
                            <option value="13-00">13:00</option>
                            <option value="17-00">17:00</option>
                            <option value="17-15">17:15</option>
                            <option value="18-00">18:00</option>
                            <option value="other">Inna</option>
                        </select>
                    </div>


                    <div class="custom-select">
                        <input type="number" name="points_value" class="input" placeholder="Wartość [pkt]" min="0" max="20" required>

                    </div>
                    <div class="custom-select">
                        <input type="date" name="date" class="input" required>
                    </div>










                    <input type="submit" value="Zapisz" name="submit">
                </form>
                <?php
                    if(isset($_POST["submit"])) {
                        switch($_POST["select_time"]) {
                            case "6-30":
                                $godzina = "6:30";
                                break;
                            case "7-00":
                                $godzina = "7:00";
                                break;
                            case "8-00":
                                $godzina = "8:00";
                                break;
                            case "9-30":
                                $godzina = "9:30";
                                break;
                            case "11-00":
                                $godzina = "11:00";
                                break;
                            case "13-00":
                                $godzina = "13:00";
                                break;
                            case "17-00":
                                $godzina = "17:00";
                                break;
                            case "17-15":
                                $godzina = "17:15";
                                break;
                            case "18-00":
                                $godzina = "18:00";
                                break;
                            case "other":
                                $godzina = "";
                                break;
                            default:
                                $godzina = "";
                                break;
                        }
                        $select_person = $_POST["select_person"];
                        $points_value = $_POST["points_value"];
                        $date = $_POST["date"];
                        $strQuery = "'$select_person','$date','$points_value','NIE','$godzina'";
                        $result = $_SESSION['conn']->query('INSERT INTO odczyt (karta,data_odczytu,ilosc_pkt,czy_odczytano,msza_godz) values ('.$strQuery.")");
                        if ($result){
                            echo "<script>alert('Dodano obecność');</script>";

                        } else {
                            echo "<script>alert('Nie udało się dodać obecności\nSpróbuj ponownie');</script>";
                        }
                        unset($_POST["submit"]);
                        unset($_POST["select_person"]);
                        unset($_POST["points_value"]);
                        unset($_POST["date"]);
                        unset($_POST["select_time"]);
                        
                    }


                ?>

            </article>

        </div>
    </main>




    <script>
        var x, i, j, l, ll, selElmnt, a, b, c;
        /*look for any elements with the class "custom-select":*/
        x = document.getElementsByClassName("custom-select");
        l = x.length;
        for (i = 0; i < l; i++) {
            selElmnt = x[i].getElementsByTagName("select")[0];
            ll = selElmnt.length;
            /*for each element, create a new DIV that will act as the selected item:*/
            a = document.createElement("DIV");
            a.setAttribute("class", "select-selected");
            a.innerHTML = selElmnt.options[selElmnt.selectedIndex].innerHTML;
            x[i].appendChild(a);
            /*for each element, create a new DIV that will contain the option list:*/
            b = document.createElement("DIV");
            b.setAttribute("class", "select-items select-hide");
            for (j = 1; j < ll; j++) {
                /*for each option in the original select element,
                create a new DIV that will act as an option item:*/
                c = document.createElement("DIV");
                c.innerHTML = selElmnt.options[j].innerHTML;
                c.addEventListener("click", function(e) {
                    /*when an item is clicked, update the original select box,
                    and the selected item:*/
                    var y, i, k, s, h, sl, yl;
                    s = this.parentNode.parentNode.getElementsByTagName("select")[0];
                    sl = s.length;
                    h = this.parentNode.previousSibling;
                    for (i = 0; i < sl; i++) {
                        if (s.options[i].innerHTML == this.innerHTML) {
                            s.selectedIndex = i;
                            h.innerHTML = this.innerHTML;
                            y = this.parentNode.getElementsByClassName("same-as-selected");
                            yl = y.length;
                            for (k = 0; k < yl; k++) {
                                y[k].removeAttribute("class");
                            }
                            this.setAttribute("class", "same-as-selected");
                            break;
                        }
                    }
                    h.click();
                });
                b.appendChild(c);
            }
            x[i].appendChild(b);
            a.addEventListener("click", function(e) {
                /*when the select box is clicked, close any other select boxes,
                and open/close the current select box:*/
                e.stopPropagation();
                closeAllSelect(this);
                this.nextSibling.classList.toggle("select-hide");
                this.classList.toggle("select-arrow-active");
            });
        }

        function closeAllSelect(elmnt) {
            /*a function that will close all select boxes in the document,
            except the current select box:*/
            var x, y, i, xl, yl, arrNo = [];
            x = document.getElementsByClassName("select-items");
            y = document.getElementsByClassName("select-selected");
            xl = x.length;
            yl = y.length;
            for (i = 0; i < yl; i++) {
                if (elmnt == y[i]) {
                    arrNo.push(i)
                } else {
                    y[i].classList.remove("select-arrow-active");
                }
            }
            for (i = 0; i < xl; i++) {
                if (arrNo.indexOf(i)) {
                    x[i].classList.add("select-hide");
                }
            }
        }
        /*if the user clicks anywhere outside the select box,
        then close all select boxes:*/
        document.addEventListener("click", closeAllSelect);
    </script>
    <script src="../button_script.js"></script>
    <script type="module" src="https://unpkg.com/ionicons@5.5.2/dist/ionicons/ionicons.esm.js"></script>
    <script nomodule src="https://unpkg.com/ionicons@5.5.2/dist/ionicons/ionicons.js"></script>
</body>

</html>